package viewmodel

type QueueCountViewModel struct {
	Count int `json:"queue_count"`
}

func NewQueueCountViewModel(count int) *QueueCountViewModel {
	return &QueueCountViewModel{Count: count}
}
