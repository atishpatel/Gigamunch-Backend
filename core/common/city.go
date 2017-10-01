package common

type City int

const (
	// Nashville is Greater Nashville Area (Nashville, Brentwood, Franklin).
	Nashville City = 0
)

func (c City) String() string {
	switch int(c) {
	case 0:
		return "Nashville"
	}
	return ""
}
