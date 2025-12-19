package songviewmodel

import "github.com/vlady-kotsev/lime-radio/transmitter/internal/domain"

type SongViewModel struct {
	Artist string `json:"artist"`
	Title  string `json:"title"`
}

func ToSongViewModel(song *domain.Song) SongViewModel {
	return SongViewModel{
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
