package domain

import "github.com/google/uuid"

type SongRequest struct {
	ID     uuid.UUID `json:"id"`
	Artist string    `json:"artist"`
	Title  string    `json:"title"`
}
