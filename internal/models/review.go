package models

import "time"

// ReviewRequest is what the client sends to POST /review
type ReviewRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
}

// Review is a saved code review (also what we return)
type Review struct {
	ID        int       `json:"id"`
	Language  string    `json:"language"`
	Code      string    `json:"code"`
	Review    string    `json:"review"`
	CreatedAt time.Time `json:"created_at"`
}
