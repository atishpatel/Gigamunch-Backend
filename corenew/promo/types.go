package promo

import (
	"time"
)

const (
	// Pending is the PromoCode state when an inquiry is pending.
	Pending State = 0
	// Used is the PromoCode state when a promo has beed used.
	Used State = 1
	// Invalid is the PromoCode state when the inquiry using the promo is invalid (not used).
	Invalid State = 2
)

// State is the state of a PromoCode.
type State int8

// Code is information for available promo codes.
type Code struct {
	Code             string    `json:"code"`
	CreatedDatetime  time.Time `json:"created_datetime"`
	FreeDelivery     bool      `json:"free_delivery"`
	PercentOff       int32     `json:"percent_off"`
	AmountOff        float32   `json:"amount_off"`
	DiscountCap      float32   `json:"discount_cap"`
	FreeDish         bool      `json:"free_dish"`
	BuyOneGetOneFree bool      `json:"buy_one_get_one_free"`
	StartDatetime    time.Time `json:"start_datetime"`
	EndDatetime      time.Time `json:"end_datetime"`
	NumUses          int32     `json:"num_uses"`
}

// UsedCode is information for a promo code usage.
type UsedCode struct {
	EaterID         string    `json:"eater_id"`
	InquiryID       int64     `json:"inquiry_id"`
	CreatedDatetime time.Time `json:"created_datetime"`
	Code            string    `json:"code"`
	State           State     `json:"state"`
}
