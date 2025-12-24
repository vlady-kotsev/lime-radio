package domain

type Song struct {
	ID     string
	Artist string
	Title  string
	Path   string
}

func NewSong(id, artist, title, path string) *Song {
	return &Song{ID: id, Artist: artist, Title: title, Path: path}
}
