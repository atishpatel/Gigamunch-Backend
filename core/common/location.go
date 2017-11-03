package common

type Location int

// TODO: move geofence to be based on ID instead of Name. Delete Nashville entry.

const (
	// Nashville is Greater Nashville Area (Nashville, Brentwood, Franklin).
	Nashville Location = 0
)

func (c Location) String() string {
	switch int(c) {
	case 0:
		return "Nashville"
	}
	return ""
}

// ID is the unique identifier for location.
func (c Location) ID() int64 {
	return int64(c)
}
