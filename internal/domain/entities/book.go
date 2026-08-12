package entities

import "time"

type Book struct {
	ID          int       `json:"id"`
	Size        int       `json:"size"`
	Title       string    `json:"title"`
	Author      string    `json:"author"`
	Status      string    `json:"status"`
	RawTextPath string    `json:"rawTextPath"`
	CreatedAt   time.Time `json:"createdAt"`
}
