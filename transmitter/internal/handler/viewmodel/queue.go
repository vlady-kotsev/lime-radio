package viewmodel

type QueueCountViewModel struct {
	Count    int `json:"queue_count"`
	MaxCount int `json:"max_queue_count"`
}

func NewQueueCountViewModel(count, maxCount int) *QueueCountViewModel {
	return &QueueCountViewModel{Count: count, MaxCount: maxCount}
}
