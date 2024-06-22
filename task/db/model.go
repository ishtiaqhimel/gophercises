package db

import "time"

type Task struct {
	ULID        string    `json:"id"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"createdAt"`
}

type CompletedTask struct {
	ULID        string    `json:"id"`
	Description string    `json:"description"`
	CompletedAt time.Time `json:"completedAt"`
}
