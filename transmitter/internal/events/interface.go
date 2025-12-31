package events

type EventSubscriberer interface {
	Start() error
	Close() error
}
