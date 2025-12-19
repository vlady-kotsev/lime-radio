package domain

type Song struct {
	Artist string
	Title  string
	Path   string
}

func NewSong(artist, title, path string) *Song {
	return &Song{Artist: artist, Title: title, Path: path}
}
