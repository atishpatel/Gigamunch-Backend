package activity

import (
	"database/sql"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/core/common"
	mysql "github.com/go-sql-driver/mysql"
)

// Activity is a subscriber activity.
type Activity struct {
	CreatedDatetime time.Time       `json:"created_datetime" db:"created_dt"`
	Date            string          `json:"date" db:"date"`
	UserID          string          `json:"user_id" db:"user_id"`
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
	Active   bool `json:"active" db:"active"`
	Skip     bool `json:"skip" db:"skip"`
	Forgiven bool `json:"forgiven" db:"forgiven"`
	// Bag detail
	ServingsNonVegetarian int8 `json:"servings_non_vegetarian" db:"servings"`
	ServingsVegetarain    int8 `json:"servings_vegetarian" db:"veg_servings"`
	ServingsChanged       bool `json:"servings_changed" db:"servings_changed"`
	First                 bool `json:"first" db:"first"`
	// Payment
	Amount             float32                `json:"amount" db:"amount"`
	AmountPaid         float32                `json:"amount_paid" db:"amount_paid"`
	DiscountAmount     float32                `json:"discount_amount" db:"discount_amount"`
	DiscountPercent    int8                   `json:"discount_percent" db:"discount_percent"`
	Paid               bool                   `json:"paid" db:"paid"`
	PaidDatetime       mysql.NullTime         `json:"paid_datetime" db:"paid_dt"`
	PaymentProvider    common.PaymentProvider `json:"payment_provider" db:"payment_provider"`
	TransactionID      string                 `json:"transaction_id" db:"transaction_id"`
	PaymentMethodToken string                 `json:"payment_method_token" db:"payment_method_token"`
	CustomerID         string                 `json:"customer_id" db:"customer_id"`
	// Refund
	Refunded            bool           `json:"refunded" db:"refunded"`
	RefundedAmount      float32        `json:"refunded_amount" db:"refunded_amount"`
	RefundedDatetime    mysql.NullTime `json:"refunded_datetime" db:"refunded_dt"`
	RefundTransactionID string         `json:"refund_transaction_id" db:"refund_transaction_id"`
	// Gift
	Gift           bool          `json:"gift" db:"gift"`
	GiftFromUserID sql.NullInt64 `json:"gift_from_user_id" db:"gift_from_user_id"`
	// Deviant
	// used for one time parties
	Deviant       bool   `json:"deviant" db:"deviant"`
	DeviantReason string `json:"deviant_reason" db:"deviant_reason"`
}

// DateParsed is Date as time.Time.
func (act *Activity) DateParsed() time.Time {
	var d time.Time
	if act.Date == "" {
		return d
	}
	d, _ = time.Parse(act.Date, DateFormat)
	return d
}
