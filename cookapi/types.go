package main

import (
	"fmt"
	"time"

	"github.com/atishpatel/Gigamunch-Backend/corenew/inquiry"
	"github.com/atishpatel/Gigamunch-Backend/corenew/item"
	"github.com/atishpatel/Gigamunch-Backend/corenew/menu"
	"github.com/atishpatel/Gigamunch-Backend/corenew/payment"
	"github.com/atishpatel/Gigamunch-Backend/errors"
)

// Inquiry is an Inquiry.
type Inquiry struct {
	inquiry.Inquiry
}

// SubMerchant is the submerchant payout info
type SubMerchant struct {
	payment.SubMerchantInfo
}

// Item is an Item.
type Item struct {
	ID int64 `json:"id,string"`
	item.Item
	NumLikes int `json:"num_likes"`
}

// Menu is a Menu.
type Menu struct {
	menu.Menu
}

// MenuWithItems is a Menu with all it's Items.
type MenuWithItems struct {
	menu.Menu
	Items []Item `json:"items"`
}

// GigatokenReq is a request with only a gigatoken input
type GigatokenReq struct {
	Gigatoken string `json:"gigatoken"`
}

func (req *GigatokenReq) gigatoken() string {
	return req.Gigatoken
}

// valid validates a req
func (req *GigatokenReq) valid() error {
	if req.Gigatoken == "" {
		return fmt.Errorf("Gigatoken is empty.")
	}
	return nil
}

// IDReq is for request with only an ID and Gigatoken.
type IDReq struct {
	ID int64 `json:"id,string"`
	GigatokenReq
}

// EmailReq is for request with only an Email and Gigatoken.
type EmailReq struct {
	Email string `json:"email"`
	GigatokenReq
}

// DateReq is for request with only a Date and Gigatoken.
type DateReq struct {
	GigatokenReq
	Date time.Time `json:"date"`
}

func (req *DateReq) valid() error {
	if req.Date.IsZero() {
		return fmt.Errorf("Date cannot be zero.")
	}
	return req.GigatokenReq.valid()
}

// ErrorOnlyResp is a response with only an error with code
type ErrorOnlyResp struct {
	Err errors.ErrorWithCode `json:"err"`
}
