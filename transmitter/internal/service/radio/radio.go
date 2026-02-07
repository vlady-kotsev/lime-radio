package radio

import (
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/go-mp3"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	BroadcastChannelLength int = 10_000
	// Target 20ms of PCM audio per chunk (will be calculated based on sample rate)
	ChunkDurationMs int = 20
	// Buffer ahead time - server stays this much ahead of playback
	BufferAheadMs int = 200
)

type Radio struct {
	logger                *zap.Logger
	config                config.RadioConfiger
	pl                    PlaylistServicer
	connections           map[uuid.UUID]*Connection
	lock                  sync.RWMutex
	currentSongSampleRate int
	group                 *errgroup.Group
}

var _ RadioServicer = (*Radio)(nil)

func NewRadio(lc fx.Lifecycle, logger *zap.Logger, config config.RadioConfiger, pl PlaylistServicer) *Radio {
	r := Radio{
		logger:      logger,
		config:      config,
		pl:          pl,
		lock:        sync.RWMutex{},
		connections: make(map[uuid.UUID]*Connection),
		group:       &errgroup.Group{},
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go func() {
				if err := r.startBroadcast(); err != nil {
					logger.Error("Broadcast failed", zap.Error(err))
				}
			}()
			return nil
		},
		OnStop: func(_ context.Context) error {
			for _, conn := range r.connections {
				r.RemoveClient(conn.ID)
			}
			return r.group.Wait()
		},
	})

	return &r
}

func (r *Radio) AddClient() *Connection {
	newConnection := NewConnection()
	r.lock.Lock()
	r.connections[newConnection.ID] = newConnection
	r.lock.Unlock()
	r.group.Go(func() error {
		return newConnection.ConnectionLoop()
	})
	return newConnection
}

func (r *Radio) RemoveClient(ID uuid.UUID) {
	r.lock.Lock()
	conn := r.connections[ID]
	close(conn.DoneChan)
	close(conn.BroadcastChannel)
	delete(r.connections, ID)
	r.lock.Unlock()
}

func (r *Radio) GetSampleRate() int {
	return r.currentSongSampleRate
}

func (r *Radio) broadcastToAllClients(data []byte) {
	r.lock.RLock()
	defer r.lock.RUnlock()

	for _, conn := range r.connections {
		select {
		case conn.BroadcastChannel <- data:
			// Successfully sent to connection's broadcast channel
		default:
			// Connection's broadcast buffer is full, skip this chunk
			r.logger.Warn("Connection broadcast buffer full, dropping audio chunk", zap.String("client_id", conn.ID.String()))
		}
	}
}

func (r *Radio) startBroadcast() error {
	songs, err := r.pl.GetAllSongs()
	if err != nil {
		return err
	}

	// Playback loop
	for {
		for _, song := range songs {
			if r.pl.GetQueueLength() > 0 {
				requestedSong, err := r.pl.DequeueSong()
				if err != nil {
					r.logger.Error("Error dequeueing song", zap.Error(err))
				} else {
					song = requestedSong
				}
			}

			f, err := os.Open(song.Path)
			if err != nil {
				r.logger.Error("Error opening song", zap.String("path", song.Path), zap.Error(err))
				continue
			}

			decoder, err := mp3.NewDecoder(f)
			if err != nil {
				r.logger.Error("Error creating decoder", zap.String("path", song.Path), zap.Error(err))
				err = f.Close()
				if err != nil {
					r.logger.Error("Error closing file", zap.String("path", f.Name()), zap.Error(err))
				}
				continue
			}

			r.currentSongSampleRate = decoder.SampleRate()
			r.logger.Info("Current Sample Rate", zap.Int("sample_rate", r.currentSongSampleRate))

			// Calculate chunk size for exactly 20ms of PCM audio
			// 16-bit stereo = 4 bytes per sample
			samplesPerChunk := (r.currentSongSampleRate * ChunkDurationMs) / 1000
			targetChunkSize := samplesPerChunk * 4

			var chunkBuffer []byte
			readBuf := make([]byte, 8192)

			// Buffer ahead calculation
			bufferAheadSamples := int64((r.currentSongSampleRate * BufferAheadMs) / 1000)
			chunksBufferedAhead := 0
			// How many 20ms chunks to buffer ahead
			targetBufferChunks := int((bufferAheadSamples * 4) / int64(targetChunkSize))

			// Reset timing per song to avoid sample rate mixing issues
			songStartTime := time.Now()
			songSamplesStreamed := int64(0)

			r.logger.Info("Starting song",
				zap.String("path", song.Path),
				zap.Int("sample_rate", r.currentSongSampleRate),
				zap.Int("chunk_size_bytes", targetChunkSize),
				zap.Int("buffer_ahead_chunks", targetBufferChunks))

			for {
				n, err := decoder.Read(readBuf)
				if err == io.EOF {
					// Process any remaining data in buffer
					if len(chunkBuffer) > 0 {
						r.broadcastToAllClients(chunkBuffer)
					}
					break
				}
				if err != nil {
					r.logger.Error("Error reading song data", zap.String("path", song.Path), zap.Error(err))
					break
				}

				if n > 0 {
					chunkBuffer = append(chunkBuffer, readBuf[:n]...)

					// Process complete chunks of exactly 20ms
					for len(chunkBuffer) >= targetChunkSize {
						chunk := make([]byte, targetChunkSize)
						copy(chunk, chunkBuffer[:targetChunkSize])

						r.broadcastToAllClients(chunk)
						samplesInThisChunk := int64(targetChunkSize / 4)
						songSamplesStreamed += samplesInThisChunk

						chunksBufferedAhead++

						// Only apply timing after buffer-ahead period
						if chunksBufferedAhead > targetBufferChunks {
							// Calculate exact expected time for this chunk
							expectedTime := songStartTime.Add(
								time.Duration((songSamplesStreamed-bufferAheadSamples)*1000/int64(r.currentSongSampleRate)) * time.Millisecond,
							)

							now := time.Now()
							if expectedTime.After(now) {
								time.Sleep(expectedTime.Sub(now))
							}
						} else if chunksBufferedAhead == targetBufferChunks {
							r.logger.Info("Buffer-ahead complete, starting timed playback",
								zap.Int("chunks_sent", chunksBufferedAhead),
								zap.Int("target_buffer_chunks", targetBufferChunks))
						}

						// Remove processed chunk from buffer
						chunkBuffer = chunkBuffer[targetChunkSize:]
					}
				}
			}
			err = f.Close()
			if err != nil {
				r.logger.Error("Error closing file", zap.String("path", f.Name()), zap.Error(err))
			}
		}
	}
}
