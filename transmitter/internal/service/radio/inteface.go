package radio

import (
	"github.com/google/uuid"
	shareddomain "github.com/vlady-kotsev/lime-radio/shared/domain"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
)

type RadioServicer interface {
	AddClient() *Connection
	RemoveClient(connectionID uuid.UUID)
	GetSampleRate() int
}

type PlaylistServicer interface {
	UpdateSongs() error
	GetAllSongs() ([]*shareddomain.Song, error)
	GetSongsPaginated(params *domain.PaginationParams) (*domain.PaginatedResult[*shareddomain.Song], error)
	GetSongsByTitleOrArtist(keyword string, params *domain.PaginationParams) (*domain.PaginatedResult[*shareddomain.Song], error)
	GetAllSongsInQueue() []*shareddomain.Song
	EnqueueSong(songID uuid.UUID) error
	DequeueSong() (*shareddomain.Song, error)
	GetQueueLength() int
}
