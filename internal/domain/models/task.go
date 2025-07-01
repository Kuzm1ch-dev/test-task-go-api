package models

type Task struct {
	ID             string      `json:id`
	ProcessingTime int64       `json:processing_time`
	Status         string      `json:status`
	Result         interface{} `json:result`
	CreatedAt      int64       `json:created_at`
	UpdatedAt      int64       `json:updated_at`
}
