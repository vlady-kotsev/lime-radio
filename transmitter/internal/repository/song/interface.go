package songrepository

import (
	"github.com/google/uuid"
	"github.com/vlady-kotsev/lime-radio/shared/domain"
)

type SongRepositorer interface {
	UpdateSongs(songs []*domain.Song) error
	GetAllSongs() ([]*SongDTO, error)
	GetSongByID(ID uuid.UUID) (*SongDTO, error)
}
