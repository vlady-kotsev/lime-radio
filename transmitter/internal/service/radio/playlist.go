package radio

import (
	"context"
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

type Playlist struct {
	songsFolder string
	songRepo    songrepository.SongRepositorer
}

var _ PlaylistServicer = (*Playlist)(nil)

func NewPlaylist(lc fx.Lifecycle, songRepo songrepository.SongRepositorer, config config.RadioConfiger) *Playlist {

	pl := Playlist{
		songsFolder: config.GetSongFolder(),
		songRepo:    songRepo,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return pl.UpdateSongs()
		},
	})
	return &pl
}

func (pl *Playlist) UpdateSongs() error {
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
		if err != nil {
			return err
		}

		title := meta.Title()
		if title == "" {
			title, _ = strings.CutSuffix(e.Name(), ".mp3")
		}
		artist := meta.Artist()
		if artist == "" {
			artist = "Unknown"
		}

		songs = append(songs, domain.NewSong(uuid.New().String(), artist, title, filePath))
	}

	return pl.songRepo.UpdateSongs(songs)
}

func (pl *Playlist) GetAllSongs() ([]*domain.Song, error) {
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
