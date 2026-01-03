package songrepository

import (
	"github.com/google/uuid"
	shareddomain "github.com/vlady-kotsev/lime-radio/shared/domain"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
)

type SongRepositorer interface {
	UpdateSongs(songs []*shareddomain.Song) error
	GetAllSongs() ([]*SongDTO, error)
	GetSongByID(ID uuid.UUID) (*SongDTO, error)
	GetSongsPaginated(params *domain.PaginationParams) ([]*SongDTO, error)
	GetSongsByTitleOrArtistPaginated(keyword string, params *domain.PaginationParams) ([]*SongDTO, error)
	CountSongs() (int, error)
	CountSongsByTitleOrArtist(keyword string) (int, error)
}
