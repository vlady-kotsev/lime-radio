package repository

import (
	_ "embed"

	"github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
)

//go:embed sql/insertSong.sql
var insertSong string

type SongRepository struct {
	storage *Storage
}

func NewSongRepository(storage *Storage) *SongRepository {
	return &SongRepository{storage: storage}
}

func (sr *SongRepository) InsertSongs(songs []*domain.Song) error {
	tx, err := sr.storage.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	stmt, err := tx.Prepare(insertSong)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, song := range songs {
		if _, err := stmt.Exec(
			song.Artist,
			song.Title,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}
