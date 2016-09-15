package menu

import (
	"fmt"
	"time"
)

const kindMenu = "Menu"

// Menu is a menu of items from a cook.
type Menu struct {
	CreatedDateTime time.Time `json:"created_datetime" datastore:",noindex"`
	EditedDateTime  time.Time `json:"edited_datetime" datastore:",noindex"`
	CookID          string    `json:"cook_id" datastore:",index"`
	Name            string    `json:"name" datastore:",noindex"`
	Color           Color     `json:"color" datastore:",noindex"`
}

type Color int32

// HexValue returns the hex value of the color
func (c Color) HexValue() string {
	val := int32(c)
	switch val {
	case 0: // green
		return "#4CAF50"
	case 1: // teal
		return "#009688"
	case 2: // cyan
		return "#00BCD4"
	}
	return "#4CAF50"
}

func (c Color) MarshalJSON() ([]byte, error) {
	v := fmt.Sprintf("\"%s\"", c.HexValue())
	return []byte(v), nil
}

func (c *Color) UnmarshalJSON(data []byte) error {
	colorString := string(data[1 : len(data)-1])
	switch colorString {
	case "teal":
		*c = 1
	case "cyan":
		*c = 2
	}
	return nil
}
