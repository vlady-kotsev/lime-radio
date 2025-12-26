package radio

import "github.com/vlady-kotsev/lime-radio/shared/domain"

type RadioServicer interface {
	AddClient() chan []byte
	RemoveClient(client chan []byte)
	GetSampleRate() int
	UpdateSongs() error
	GetAllSongs() ([]*domain.Song, error)
}

type PlaylistServicer interface {
	UpdateSongs() error
	GetAllSongs() ([]*domain.Song, error)
}
