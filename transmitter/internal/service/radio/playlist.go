package radio

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"github.com/google/uuid"
	"github.com/vlady-kotsev/lime-radio/shared/domain"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	songrepository "github.com/vlady-kotsev/lime-radio/transmitter/internal/repository/song"
	"go.uber.org/fx"
)

type PlaylistService struct {
	config      config.RadioConfiger
	songsFolder string
	songRepo    songrepository.SongRepositorer
	songQueue   []*domain.Song
}

var _ PlaylistServicer = (*PlaylistService)(nil)

func NewPlaylist(lc fx.Lifecycle, songRepo songrepository.SongRepositorer, config config.RadioConfiger) *PlaylistService {

	pl := PlaylistService{
		songsFolder: config.GetSongFolder(),
		songRepo:    songRepo,
		config:      config,
		songQueue:   make([]*domain.Song, 0, config.GetQueueMaxCount()),
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return pl.UpdateSongs()
		},
	})
	return &pl
}

func (pl *PlaylistService) UpdateSongs() error {
	entries, err := os.ReadDir(pl.songsFolder)
	if err != nil {
		return err
	}
	var songs []*domain.Song
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		filePath := filepath.Join(pl.songsFolder, e.Name())

		cleanPath := filepath.Clean(filePath)
		if !strings.HasPrefix(cleanPath, filepath.Clean(pl.songsFolder)) {
			continue
		}

		f, err := os.Open(cleanPath)
		if err != nil {
			return err
		}
		defer f.Close()

		meta, err := tag.ReadFrom(f)

		var title, artist string
		if err != nil {
			title, _ = strings.CutSuffix(e.Name(), ".mp3")
			artist = "Unknown"
		} else {
			title = meta.Title()
			if title == "" {
				title, _ = strings.CutSuffix(e.Name(), ".mp3")
			}
			artist = meta.Artist()
			if artist == "" {
				artist = "Unknown"
			}
		}

		songs = append(songs, domain.NewSong(uuid.New().String(), artist, title, filePath))
	}

	return pl.songRepo.UpdateSongs(songs)
}

func (pl *PlaylistService) GetAllSongs() ([]*domain.Song, error) {
	songDTOs, err := pl.songRepo.GetAllSongs()
	if err != nil {
		return []*domain.Song{}, nil
	}

	var songs []*domain.Song
	for _, dto := range songDTOs {
		songs = append(songs, dto.ToDomain())
	}

	return songs, nil
}

func (pl *PlaylistService) EnqueueSong(songID uuid.UUID) error {
	songDTO, err := pl.songRepo.GetSongByID(songID)
	if err != nil {
		return err
	}

	pl.songQueue = append(pl.songQueue, songDTO.ToDomain())
	return nil
}

func (pl *PlaylistService) DequeueSong() (*domain.Song, error) {
	if len(pl.songQueue) == 0 {
		return nil, fmt.Errorf("queue is empty")
	}
	song := pl.songQueue[0]
	pl.songQueue = pl.songQueue[1:]

	return song, nil
}

func (pl *PlaylistService) GetQueueLength() int {
	return len(pl.songQueue)
}

func (pl *PlaylistService) GetAllSongsInQueue() []*domain.Song {
	return pl.songQueue
}
