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
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
	"go.uber.org/fx"
	"go.uber.org/zap"
)

type Radio struct {
	config     *config.Config
	clients    map[chan []byte]bool
	mutex      sync.RWMutex
	sampleRate int
	logger     *zap.Logger
	songs      []*domain.Song
	pl         *Playlist
}

func NewStation(lc fx.Lifecycle, logger *zap.Logger, pl *Playlist, config *config.Config) (*Radio, error) {
	songs, err := pl.getAllSongs()
	if err != nil {
		return nil, err
	}
	r := &Radio{
		config:  config,
		clients: make(map[chan []byte]bool),
		logger:  logger,
		songs:   songs,
		pl:      pl,
	}

	lc.Append(fx.Hook{
		OnStart: func(_ context.Context) error {
			go r.startBroadcast()
			return nil
		},
	})

	return r, nil
}

func (r *Radio) AddClient() chan []byte {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	client := make(chan []byte, 10240)
	r.clients[client] = true
	r.logger.Info("Client connected", zap.Int("total_clients", len(r.clients)))
	return client
}

func (r *Radio) RemoveClient(client chan []byte) {
	r.mutex.Lock()
	defer r.mutex.Unlock()

	if _, exists := r.clients[client]; exists {
		delete(r.clients, client)
		close(client)
		r.logger.Info("Client disconnected", zap.Int("total_clients", len(r.clients)))
	}
}

func (r *Radio) broadcast(data []byte) {
	r.mutex.RLock()
	defer r.mutex.RUnlock()

	dataCopy := make([]byte, len(data))
	copy(dataCopy, data)

	for client := range r.clients {
		select {
		case client <- dataCopy:
		default:
			// Client buffer full - normal for slow clients
		}
	}
}

func (r *Radio) startBroadcast() {

	if len(r.songs) == 0 {
		r.logger.Fatal("No MP3 files found in songs directory")
		return
	}

	for {
		for _, song := range r.songs {
			songName := filepath.Base(song.Path)
			songName = strings.TrimSuffix(songName, ".mp3")
			r.logger.Info("Now playing", zap.String("song", songName))

			f, err := os.Open(song.Path)
			if err != nil {
				r.logger.Error("Error opening song", zap.String("path", song.Path), zap.Error(err))
				continue
			}

			decoder, err := mp3.NewDecoder(f)
			if err != nil {
				r.logger.Error("Error creating decoder", zap.String("path", song.Path), zap.Error(err))
				f.Close()
				continue
			}

			r.sampleRate = decoder.SampleRate()

			buf := make([]byte, 1024)

			intervalMs := calculateStreamInterval(decoder.SampleRate())
			ticker := time.NewTicker(time.Duration(intervalMs) * time.Millisecond)

			for range ticker.C {
				n, err := decoder.Read(buf)
				if err == io.EOF {
					ticker.Stop()
					break
				}
				if err != nil {
					r.logger.Error("Error reading song data", zap.String("path", song.Path), zap.Error(err))
					ticker.Stop()
					break
				}
				if n > 0 {
					r.broadcast(buf[:n])
				}
			}
			f.Close()
		}
	}
}

func (r *Radio) GetSampleRate() int {
	return r.sampleRate
}

func (r *Radio) UpdateSongs() error {
	err := r.pl.updateSongs()
	if err != nil {
		return err
	}
	songs, err := r.pl.getAllSongs()
	if err != nil {
		return err
	}

	r.songs = songs
	return nil
}

func (r *Radio) GetAllSongs() ([]*domain.Song, error) {
	return r.pl.getAllSongs()
}
