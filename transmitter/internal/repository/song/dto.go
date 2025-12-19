package songrepository

import "github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"

type SongDTO struct {
	Id     uint   `db:"id"`
	Artist string `db:"artist"`
	Title  string `db:"title"`
	Path   string `db:"path"`
}

func (dto *SongDTO) ToDomain() *domain.Song {
	return domain.NewSong(dto.Artist, dto.Title, dto.Path)
}

func FromDomain(song *domain.Song) *SongDTO {
	return &SongDTO{
		Artist: song.Artist,
		Title:  song.Title,
		Path:   song.Path,
	}
}
