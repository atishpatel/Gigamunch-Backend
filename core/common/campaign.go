package common

import "time"

// Campaign is a campaign a subscriber was a part of.
type Campaign struct {
	Source    string    `json:"source"`
	Medium    string    `json:"medium"`
	Campaign  string    `json:"campaign"`
	Term      string    `json:"term" datastore:",noindex"`
	Content   string    `json:"content" datastore:",noindex"`
	Timestamp time.Time `json:"timestamp"`
}
