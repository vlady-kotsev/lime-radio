package radio

import (
	"github.com/google/uuid"
	"github.com/vlady-kotsev/lime-radio/shared/domain"
)

type RadioServicer interface {
	AddClient() *Connection
	RemoveClient(connectionID uuid.UUID)
	GetSampleRate() int
}

type PlaylistServicer interface {
	UpdateSongs() error
	GetAllSongs() ([]*domain.Song, error)
	GetAllSongsByTitleOrArtist(keyword string) ([]*domain.Song, error)
	GetAllSongsInQueue() []*domain.Song
	EnqueueSong(songID uuid.UUID) error
	DequeueSong() (*domain.Song, error)
	GetQueueLength() int
}
