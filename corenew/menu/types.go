package menu

import (
	"fmt"
	"math/rand"
	"time"
)

const kindMenu = "Menu"

var colors = []string{
	"#66BB6A", // Green      400
	"#689F38", // LightGreen 700
	"#EA9E22", // Yellow     850
	"#F57C00", // Orange     700
	"#FF7D54", // DeepOrange 350
	"#EF5350", // Red        400
	"#EA4F83", // Pink       350
	"#B159BF", // Purple     350
	"#9575CD", // DeepPurple 300
	"#5C6BC0", // Indigo     400
	"#0192D9", // Light Blue 650
	"#00ACC1", // Cyan       600
	"#26A69A", // Teal       400
}

// Menu is a menu of items from a cook.
type Menu struct {
	ID              int64     `json:"id,string" datastore:",noindex"`
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
	return Color(rand.Int31n(int32(len(colors))) + 1)
}

func (c Color) isZero() bool {
	return int32(c) == 0
}

// HexValue returns the hex value of the color
func (c Color) HexValue() string {
	val := int(c)
	if val <= 0 || val >= len(colors) {
		return colors[0]
	}
	return colors[val-1]
}

// MarshalJSON converts Color into a HexValue
func (c Color) MarshalJSON() ([]byte, error) {
	v := fmt.Sprintf("\"%s\"", c.HexValue())
	return []byte(v), nil
}

// UnmarshalJSON converts Color into a int32
func (c *Color) UnmarshalJSON(data []byte) error {
	colorString := string(data[1 : len(data)-1])
	for i, color := range colors {
		if color == colorString {
			*c = Color(i + 1)
			return nil
		}
	}
	*c = NewColor()
	return nil
}
