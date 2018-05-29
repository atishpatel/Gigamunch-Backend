package deliveries

import "time"

// Delivery is a delivery for a subscriber.
type Delivery struct {
	Date            string    `json:"date,omitempty" db:"date"`
	SubEmail        string    `json:"sub_email,omitempty" db:"sub_email"`
	CreatedDatetime time.Time `json:"created_datetime,omitempty" db:"created_dt"`
	UpdatedDatetime time.Time `json:"updated_datetime,omitempty" db:"updated_dt"`
	DriverID        int64     `json:"driver_id,omitempty" db:"driver_id"`
	DriverName      string    `json:"driver_name,omitempty" db:"driver_name"`
	DriverEmail     string    `json:"driver_email,omitempty" db:"driver_email"`
	SubID           int64     `json:"sub_id,omitempty" db:"sub_id"`
	Order           int       `json:"order,omitempty" db:"delivery_order"`
	Success         bool      `json:"success,omitempty" db:"success"`
	Fail            bool      `json:"fail,omitempty" db:"fail"`
}
