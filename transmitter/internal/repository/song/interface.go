package songrepository

import "github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"

type SongRepositorer interface {
	UpdateSongs(songs []*domain.Song) error
	GetAllSongs() ([]*SongDTO, error)
}
