package songrepository

import "github.com/vlady-kotsev/lime-radio/shared/domain"

type SongRepositorer interface {
	UpdateSongs(songs []*domain.Song) error
	GetAllSongs() ([]*SongDTO, error)
}
