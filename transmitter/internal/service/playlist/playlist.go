package playlist

import (
	"context"
	"log"
	"os"
	"path/filepath"

	"github.com/dhowden/tag"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/repository"
	"go.uber.org/fx"
)

type Playlist struct {
	repo *repository.SongRepository
}

func NewPlaylist(lc fx.Lifecycle, songRepo *repository.SongRepository) *Playlist {
	songsFolder := "songs"
	pl := Playlist{
		repo: songRepo,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return pl.updateSongs(songsFolder)
		},
	})
	return &pl
}

func (pl *Playlist) updateSongs(songsFolder string) error {
	entries, err := os.ReadDir(songsFolder)
	if err != nil {
		log.Fatal(err)
	}
	var songs []*domain.Song
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		f, err := os.Open(filepath.Join(songsFolder, e.Name()))
		if err != nil {
			return err
		}
		defer f.Close()

		meta, err := tag.ReadFrom(f)
		if err != nil {
			return err
		}
		songs = append(songs, domain.NewSong(meta.Artist(), meta.Title()))
	}

	return pl.repo.InsertSongs(songs)
}

func (pl *Playlist) GetPlaylist() ([]string, error) {
	return filepath.Glob("songs/*.mp3")
}
