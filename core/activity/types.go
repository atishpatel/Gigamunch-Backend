package activity

import (
	"database/sql"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	mysql "github.com/go-sql-driver/mysql"
)

type Activity struct {
	CreatedDatetime time.Time       `json:"created_datetime" db:"created_dt"`
	Date            string          `json:"date" db:"date"`
	UserID          int64           `json:"user_id" db:"user_id"`
	Email           string          `json:"email" db:"email"`
	FirstName       string          `json:"first_name" db:"first_name"`
	LastName        string          `json:"last_name" db:"last_name"`
	Location        common.Location `json:"location" db:"location"`
	// Address
	AddressChanged bool    `json:"address_changed" db:"addr_changed"`
	AddressAPT     string  `json:"address_apt" db:"addr_apt"`
	AddressString  string  `json:"address_string" db:"addr_string"`
	Zip            string  `json:"zip" db:"zip"`
	Latitude       float64 `json:"latitude,string" db:"lat"`
	Longitude      float64 `json:"longitude,string" db:"long"`
	// Detail
	Active bool `json:"active" db:"active"`
	Skip   bool `json:"skip" db:"skip"`
	// Bag detail
	Servings          int8 `json:"servings" db:"servings"`
	VegetrainServings bool `json:"vegetarian_servings" db:"veg_servings"`
	ServingsChanged   int8 `json:"servings_changed" db:"servings_changed"`
	First             bool `json:"first" db:"first"`
	// Payment
	Amount         float32        `json:"amount" db:"amount"`
	AmountPaid     float32        `json:"amount_paid" db:"amount_paid"`
	DiscountAmount float32        `json:"discount_amount" db:"discount_amount"`
	Paid           bool           `json:"paid" db:"paid"`
	PaidDatetime   mysql.NullTime `json:"paid_datetime" db:"paid_dt"`
	TransactionID  string         `json:"transaction_id" db:"transaction_id"`
	// Refund
	Refunded            bool           `json:"refunded" db:"refunded"`
	RefundedAmount      float32        `json:"refunded_amount" db:"refunded_amount"`
	RefundedDatetime    mysql.NullTime `json:"refunded_datetime" db:"refunded_dt"`
	RefundTransactionID string         `json:"refund_transaction_id" db:"refund_transaction_id"`
	// CouponID            sql.NullInt64          `json:"coupon_id" db:"coupon_id"`
	PaymentProvider common.PaymentProvider `json:"payment_provider" db:"payment_provider"`
	Forgiven        bool                   `json:"forgiven" db:"forgiven"`
	// Gift
	Gift           bool          `json:"gift" db:"gift"`
	GiftFromUserID sql.NullInt64 `json:"gift_from_user_id" db:"gift_from_user_id"`
	// Deviant
	// used for one time parties
	Deviant       bool   `json:"deviant" db:"deviant"`
	DeviantReason string `json:"deviant_reason" db:"deviant_reason"`
	// Driver
	AssignedDriverID sql.NullInt64  `json:"assigned_driver_id" db:"assigned_driver_id"`
	Delivered        bool           `json:"delivered" db:"delivered"`
	DeliveryDatetime mysql.NullTime `json:"delivery_datetime" db:"delivery_dt"`
}
