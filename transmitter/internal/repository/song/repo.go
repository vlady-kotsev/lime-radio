package songrepository

import (
	_ "embed"
	"fmt"

	"github.com/google/uuid"
	"github.com/jmoiron/sqlx"
	"github.com/vlady-kotsev/lime-radio/shared/domain"
	transmitterdomain "github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"
	"github.com/vlady-kotsev/lime-radio/transmitter/internal/repository"
)

var (
	//go:embed sql/insertSong.sql
	insertSong string
	//go:embed sql/getSongs.sql
	getSongs string
	//go:embed sql/deleteSong.sql
	deleteSong string
	//go:embed sql/getSongById.sql
	getSongById string
	//go:embed sql/getSongsPaginated.sql
	getSongsPaginated string
	//go:embed sql/getSongsByTitleOrArtistPaginated.sql
	getSongsByTitleOrArtistPaginated string
	//go:embed sql/countSongs.sql
	countSongs string
	//go:embed sql/countSongsByTitleOrArtist.sql
	countSongsByTitleOrArtist string
)

type SongMapEntry struct {
	Artist string
	Title  string
}

type SongRepository struct {
	storage repository.Storager
}

var _ SongRepositorer = (*SongRepository)(nil)

func NewSongRepository(storage repository.Storager) *SongRepository {
	return &SongRepository{storage: storage}
}

func (sr *SongRepository) insertSongs(tx *sqlx.Tx, dtos []*SongDTO) error {
	stmt, err := tx.Prepare(insertSong)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, dto := range dtos {
		if dto.Id == "" {
			dto.Id = uuid.New().String()
		}
		if _, err := stmt.Exec(
			dto.Id,
			dto.Artist,
			dto.Title,
			dto.Path,
		); err != nil {
			return err
		}
	}
	return nil
}

func (sr *SongRepository) deleteSongs(tx *sqlx.Tx, dtos []*SongDTO) error {
	stmt, err := tx.Prepare(deleteSong)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, dto := range dtos {
		if _, err := stmt.Exec(
			dto.Artist,
			dto.Title,
		); err != nil {
			return err
		}
	}
	return nil
}

func (sr *SongRepository) GetSongByID(ID uuid.UUID) (*SongDTO, error) {
	var dto SongDTO
	err := sr.storage.Get(&dto, getSongById, ID.String())
	if err != nil {
		return nil, err
	}
	return &dto, nil
}

func (sr *SongRepository) UpdateSongs(songs []*domain.Song) error {
	songsMap := make(map[SongMapEntry]SongDTO)
	for _, song := range songs {
		songsMap[SongMapEntry{
			Title:  song.Title,
			Artist: song.Artist,
		}] = SongDTO{Artist: song.Artist, Title: song.Title, Path: song.Path}
	}

	tx, err := sr.storage.Beginx()
	if err != nil {
		return err
	}
	defer func() {
		_ = tx.Rollback()
	}()

	var dtos []SongDTO
	err = tx.Select(&dtos, getSongs)
	if err != nil {
		return err
	}
	dtosMap := make(map[SongMapEntry]SongDTO)
	for _, dto := range dtos {
		dtosMap[SongMapEntry{
			Title:  dto.Title,
			Artist: dto.Artist,
		}] = dto
	}

	deletedSongs := make([]*SongDTO, 0, len(dtosMap))
	addedSongs := make([]*SongDTO, 0, len(songsMap))
	for songEntry, songDTO := range songsMap {
		if _, ok := dtosMap[songEntry]; !ok {
			addedSongs = append(addedSongs, &songDTO)
		}
	}
	for songEntry, songDTO := range dtosMap {
		if _, ok := songsMap[songEntry]; !ok {
			deletedSongs = append(deletedSongs, &songDTO)
		}
	}

	err = sr.insertSongs(tx, addedSongs)
	if err != nil {
		return err
	}

	err = sr.deleteSongs(tx, deletedSongs)
	if err != nil {
		return err
	}

	return tx.Commit()
}

func (sr *SongRepository) GetAllSongs() ([]*SongDTO, error) {
	var dtos []*SongDTO
	if err := sr.storage.Select(&dtos, getSongs); err != nil {
		return nil, err
	}

	return dtos, nil
}

func (sr *SongRepository) GetSongsPaginated(params *transmitterdomain.PaginationParams) ([]*SongDTO, error) {
	var dtos []*SongDTO
	if err := sr.storage.Select(&dtos, getSongsPaginated, params.GetLimit(), params.GetOffset()); err != nil {
		return nil, err
	}
	return dtos, nil
}

func (sr *SongRepository) GetSongsByTitleOrArtistPaginated(keyword string, params *transmitterdomain.PaginationParams) ([]*SongDTO, error) {
	wildcard := fmt.Sprintf("%%%s%%", keyword)
	var dtos []*SongDTO
	if err := sr.storage.Select(&dtos, getSongsByTitleOrArtistPaginated, wildcard, wildcard, params.GetLimit(), params.GetOffset()); err != nil {
		return nil, err
	}
	return dtos, nil
}

func (sr *SongRepository) CountSongs() (int, error) {
	var count int
	if err := sr.storage.Get(&count, countSongs); err != nil {
		return 0, err
	}
	return count, nil
}

func (sr *SongRepository) CountSongsByTitleOrArtist(keyword string) (int, error) {
	wildcard := fmt.Sprintf("%%%s%%", keyword)
	var count int
	if err := sr.storage.Get(&count, countSongsByTitleOrArtist, wildcard, wildcard); err != nil {
		return 0, err
	}
	return count, nil
}
