package gigamuncher

import "github.com/atishpatel/Gigamunch-Backend/core/review"

/*
 * This file is for types shared between multiple files.
 */

// Review is a review
type Review struct {
	ID            int `json:"id"`
	review.Review     // embedded
}
