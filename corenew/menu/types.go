package menu

import (
	"fmt"
	"math/rand"
	"time"
)

const kindMenu = "Menu"

// Menu is a menu of items from a cook.
type Menu struct {
	CreatedDateTime time.Time `json:"created_datetime" datastore:",noindex"`
	// EditedDateTime  time.Time `json:"edited_datetime" datastore:",noindex"`
	CookID string `json:"cook_id" datastore:",index"`
	Name   string `json:"name" datastore:",noindex"`
	Color  Color  `json:"color" datastore:",noindex"`
}

// Color is a color for the menu
type Color int32

// NewColor returns a new random color
func NewColor() Color {
	return Color(rand.Int31n(7) + 1)
}

func (c Color) isZero() bool {
	return int32(c) == 0
}

// HexValue returns the hex value of the color
func (c Color) HexValue() string {
	val := int32(c)
	switch val {
	case 1: // green
		return "#4CAF50"
	case 2: // teal
		return "#009688"
	case 3: // cyan
		return "#00BCD4"
	}
	return "#4CAF50"
}

// MarshalJSON converts Color into a HexValue
func (c Color) MarshalJSON() ([]byte, error) {
	v := fmt.Sprintf("\"%s\"", c.HexValue())
	return []byte(v), nil
}

// UnmarshalJSON converts Color into a int32
func (c *Color) UnmarshalJSON(data []byte) error {
	colorString := string(data[1 : len(data)-1])
	switch colorString {
	case "#009688":
		*c = 2
	case "#00BCD4":
		*c = 3
	default:
		*c = NewColor()
	}
	return nil
}
