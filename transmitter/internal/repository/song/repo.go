package songrepository

import (
	_ "embed"

	"github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/repository"
)

var (
	//go:embed sql/insertSong.sql
	insertSong string
	//go:embed sql/getSongs.sql
	getSongs string
)

type SongRepository struct {
	storage *repository.Storage
}

func NewSongRepository(storage *repository.Storage) *SongRepository {
	return &SongRepository{storage: storage}
}

func (sr *SongRepository) InsertSongs(songs []*domain.Song) error {
	tx, err := sr.storage.DB.Begin()
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
			song.Path,
		); err != nil {
			return err
		}
	}
	return tx.Commit()
}

func (sr *SongRepository) GetAllSongs() ([]*SongDTO, error) {
	var dtos []*SongDTO
	if err := sr.storage.DB.Select(&dtos, getSongs); err != nil {
		return nil, err
	}

	return dtos, nil
}
