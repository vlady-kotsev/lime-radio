package radio

import (
	"context"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/dhowden/tag"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/config"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
	songrepository "github.com/vlady-kotsev/lime-radio/transmitter/internal/repository/song"
	"go.uber.org/fx"
)

type Playlist struct {
	songsFolder string
	songRepo    *songrepository.SongRepository
}

func NewPlaylist(lc fx.Lifecycle, songRepo *songrepository.SongRepository, config *config.Config) *Playlist {

	pl := Playlist{
		songsFolder: config.App.SongsFolder,
		songRepo:    songRepo,
	}

	lc.Append(fx.Hook{
		OnStart: func(context.Context) error {
			return pl.updateSongs()
		},
	})
	return &pl
}

func (pl *Playlist) updateSongs() error {
	entries, err := os.ReadDir(pl.songsFolder)
	if err != nil {
		log.Fatal(err)
	}
	var songs []*domain.Song
	for _, e := range entries {
		if e.IsDir() {
			continue
		}
		filePath := filepath.Join(pl.songsFolder, e.Name())
		f, err := os.Open(filePath)
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

		songs = append(songs, domain.NewSong(artist, title, filePath))
	}

	return pl.songRepo.InsertSongs(songs)
}

func (pl *Playlist) getAllSongs() ([]*domain.Song, error) {
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
