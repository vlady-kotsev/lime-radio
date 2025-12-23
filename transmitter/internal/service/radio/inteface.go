package radio

import "github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"

type Radioer interface {
	AddClient() chan []byte
	RemoveClient(client chan []byte)
	GetSampleRate() int
	UpdateSongs() error
	GetAllSongs() ([]*domain.Song, error)
}

type Playlister interface {
	UpdateSongs() error
	GetAllSongs() ([]*domain.Song, error)
}
