package viewmodel

import "github.com/vlady-kotsev/lime-radio/shared/domain"

type SongViewModel struct {
	ID     string `json:"id"`
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

func ToSongViewModel(song *domain.Song) SongViewModel {
	return SongViewModel{
		ID:     song.ID,
		Artist: song.Artist,
		Title:  song.Title,
	}
}

func ToSongViewModels(songs []*domain.Song) []SongViewModel {
	viewModels := make([]SongViewModel, len(songs))
	for i, song := range songs {
		viewModels[i] = ToSongViewModel(song)
	}
	return viewModels
}
