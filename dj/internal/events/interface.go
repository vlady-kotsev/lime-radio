package events

type EventPublisherer interface {
	PublishMessage(payload string) error
}
