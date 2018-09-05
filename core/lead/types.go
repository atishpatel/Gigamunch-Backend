package lead

import "time"

const kind = "Lead"

// Lead is a customer lead
type Lead struct {
	CreatedDatetime time.Time
	Email           string
}
