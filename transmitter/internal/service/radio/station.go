package radio

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/hajimehoshi/go-mp3"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/service/playlist"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Station struct {
	clients    map[chan []byte]bool
	mutex      sync.RWMutex
	sampleRate int
	logger     *zap.Logger
}

func NewStation(lc fx.Lifecycle, logger *zap.Logger, pl *playlist.Playlist) *Station {
	s := &Station{
		clients: make(map[chan []byte]bool),
		logger:  logger,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go s.startBroadcast(pl)
			return nil
		},
	})

	return s
}

func (s *Station) AddClient() chan []byte {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	client := make(chan []byte, 100000)
	s.clients[client] = true
	s.logger.Info("Client connected", zap.Int("total_clients", len(s.clients)))
	return client
}

func (s *Station) RemoveClient(client chan []byte) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.clients[client]; exists {
		delete(s.clients, client)
		close(client)
		s.logger.Info("Client disconnected", zap.Int("total_clients", len(s.clients)))
	}
}

func (s *Station) broadcast(data []byte) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	for client := range s.clients {
		select {
		case client <- dataCopy:
		default:
			// Client buffer full - normal for slow clients
		}
	}
}

func (s *Station) startBroadcast(pl *playlist.Playlist) {
	playlist, err := pl.GetPlaylist()
	if err != nil || len(playlist) == 0 {
		s.logger.Fatal("No MP3 files found in songs directory")
		return
	}

	for {
		for _, songPath := range playlist {
			songName := filepath.Base(songPath)
			songName = strings.TrimSuffix(songName, ".mp3")
			s.logger.Info("Now playing", zap.String("song", songName))

			f, err := os.Open(songPath)
			if err != nil {
				s.logger.Error("Error opening song", zap.String("path", songPath), zap.Error(err))
				continue
			}

			decoder, err := mp3.NewDecoder(f)
			if err != nil {
				s.logger.Error("Error creating decoder", zap.String("path", songPath), zap.Error(err))
				f.Close()
				continue
			}

			s.sampleRate = decoder.SampleRate()

			buf := make([]byte, 1024)
			ticker := time.NewTicker(time.Millisecond * 3)

			for range ticker.C {
				n, err := decoder.Read(buf)
				if err == io.EOF {
					ticker.Stop()
					break
				}
				if err != nil {
					s.logger.Error("Error reading song data", zap.String("path", songPath), zap.Error(err))
					ticker.Stop()
					break
				}
				if n > 0 {
					s.broadcast(buf[:n])
				}
			}
			f.Close()
		}
	}
}

func (s *Station) GetSampleRate() int {
	return s.sampleRate
}

func (s *Station) StreamToClient() chan []byte {
	return s.AddClient()
}
