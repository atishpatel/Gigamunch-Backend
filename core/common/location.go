package common

type Location int

const (
	// Nashville is Greater Nashville Area (Nashville, Brentwood, Franklin).
	Nashville Location = 1
)

func (c Location) String() string {
	switch int64(c) {
	case int64(Nashville):
		return "Nashville"
	}
	return ""
}

// ID is the unique identifier for location.
func (c Location) ID() string {
	switch int64(c) {
	case int64(Nashville):
		return "Nashville"
	}
	return ""
}
