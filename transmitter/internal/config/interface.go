package config

type RadioConfiger interface {
	GetSongFolder() string
	GetQueueMaxCount() int
}
