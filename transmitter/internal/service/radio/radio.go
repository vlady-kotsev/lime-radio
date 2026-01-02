package radio

import (
	"bytes"
	"context"
	"io"
	"os"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hajimehoshi/go-mp3"
	"github.com/vlady-kotsev/lime-radio/shared/domain"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	"go.uber.org/fx"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
)

const (
	SongQueueLength        int = 10
	BroadcastChannelLength int = 10_000
	ChunkSize              int = 1024
)

type Radio struct {
	logger                *zap.Logger
	config                config.RadioConfiger
	pl                    PlaylistServicer
	connections           map[uuid.UUID]*Connection
	lock                  sync.RWMutex
	currentSongSampleRate int
	songQueue             chan *domain.Song
	broadcastChan         chan []byte
	group                 *errgroup.Group
}

var _ RadioServicer = (*Radio)(nil)

func NewRadio(lc fx.Lifecycle, logger *zap.Logger, config config.RadioConfiger, pl PlaylistServicer) *Radio {
	r := Radio{
		logger:        logger,
		config:        config,
		pl:            pl,
		lock:          sync.RWMutex{},
		connections:   make(map[uuid.UUID]*Connection),
		songQueue:     make(chan *domain.Song, SongQueueLength),
		broadcastChan: make(chan []byte, BroadcastChannelLength),
		group:         &errgroup.Group{},
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
			close(r.broadcastChan)
			return r.group.Wait()
		},
	})

	return &r
}

func (r *Radio) AddClient() *Connection {
	newConnection := NewConnection(r.broadcastChan)
	r.lock.Lock()
	r.connections[newConnection.ID] = newConnection
	r.lock.Unlock()
	r.group.Go(newConnection.ConnectionLoop)
	return newConnection
}

func (r *Radio) RemoveClient(ID uuid.UUID) {
	r.lock.Lock()
	conn := r.connections[ID]
	close(conn.DoneChan)
	delete(r.connections, ID)
	r.lock.Unlock()
}

func (r *Radio) GetSampleRate() int {
	return r.currentSongSampleRate
}

func (r *Radio) startBroadcast() error {
	songs, err := r.pl.GetAllSongs()
	if err != nil {
		return err
	}
	// Playback loop
	for {
		for _, song := range songs {
			// Check if we have songs in queue
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
			defer f.Close()

			data, err := io.ReadAll(f)
			if err != nil {
				r.logger.Error("Error reading song file", zap.String("path", song.Path), zap.Error(err))
				continue
			}

			reader := bytes.NewReader(data)
			decoder, err := mp3.NewDecoder(reader)
			if err != nil {
				r.logger.Error("Error creating decoder", zap.String("path", song.Path), zap.Error(err))
				continue
			}
			r.currentSongSampleRate = decoder.SampleRate()

			buf := make([]byte, ChunkSize)
			for {
				bytesRead, err := decoder.Read(buf)
				if err == io.EOF {
					break
				}
				if err != nil {
					r.logger.Error("Error reading song data", zap.String("path", song.Path), zap.Error(err))
					break
				}

				if bytesRead > 0 {
					dataCopy := make([]byte, bytesRead)
					copy(dataCopy, buf[:bytesRead])

					r.broadcastChan <- dataCopy

					waitTime := calculateStreamIntervalForBytes(decoder.SampleRate(), bytesRead)
					time.Sleep(waitTime)
				}
			}
		}
	}
}
