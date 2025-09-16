package models

import "time"

type URL struct {
	ID        int       `json:"id"`
	Original  string    `json:"original"`
	Short     string    `json:"short"`
	Clicks    int       `json:"clicks"`
	CreatedAt time.Time `json:"created_at"`
}
