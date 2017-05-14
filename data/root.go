package data

import "time"

// ActivityRequest represents an API request for activity
type ActivityRequest struct {
	StartTime time.Time `json:"starttime"`
	EndTime   time.Time `json:"endtime"`
}

// ActivityResponse represents an API response
type ActivityResponse struct {
	ID      int    `json:"id"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}
