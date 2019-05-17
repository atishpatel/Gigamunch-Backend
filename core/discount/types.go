package discount

import (
	"time"
)

// Discount is a subscriber Discount.
type Discount struct {
	ID              int64     `json:"id" db:"id"`
	CreatedDatetime time.Time `json:"created_datetime,string" db:"created_dt"`
	UserID          string    `json:"user_id" db:"user_id"`
	Email           string    `json:"email" db:"email"`
	FirstName       string    `json:"first_name" db:"first_name"`
	LastName        string    `json:"last_name" db:"last_name"`
	DateUsed        string    `json:"date_used" db:"date_used"`
	DiscountAmount  float32   `json:"discount_amount" db:"discount_amount"`
	DiscountPercent int8      `json:"discount_percent" db:"discount_percent"`
}

// IsUsed returns if discount is used.
func (d *Discount) IsUsed() bool {
	return d.DateUsed[:10] > "1001-01-01"
}
