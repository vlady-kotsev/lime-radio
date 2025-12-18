package domain

type Song struct {
	Artist string
	Title  string
}

func NewSong(artist, title string) *Song {
	return &Song{Artist: artist, Title: title}
}
