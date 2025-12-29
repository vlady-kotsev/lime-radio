package radio

import (
	"github.com/google/uuid"
	"github.com/vlady-kotsev/lime-radio/shared/domain"
)

type RadioServicer interface {
	AddClient() *Connection
	RemoveClient(connectionID uuid.UUID)
	GetSampleRate() int
	UpdateSongs() error
	GetAllSongs() ([]*domain.Song, error)
}

type PlaylistServicer interface {
	UpdateSongs() error
	GetAllSongs() ([]*domain.Song, error)
}
