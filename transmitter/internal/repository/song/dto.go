package songrepository

import "github.com/vlady-kotsev/lime-radio/shared/domain"

type SongDTO struct {
	Id     string `db:"id"`
	Artist string `db:"artist"`
	Title  string `db:"title"`
	Path   string `db:"path"`
}

func (dto *SongDTO) ToDomain() *domain.Song {
	return domain.NewSong(dto.Id, dto.Artist, dto.Title, dto.Path)
}

func FromDomain(song *domain.Song) *SongDTO {
	return &SongDTO{
		Id:     song.ID,
		Artist: song.Artist,
		Title:  song.Title,
		Path:   song.Path,
	}
}
