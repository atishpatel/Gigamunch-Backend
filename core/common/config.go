package common

const (
	prod  = "gigamunch-omninexus"
	stage = "gigamunch-omninexus-dev"
)

// IsDev returns true if it's dev env.
func IsDev(projID string) bool {
	return projID != prod && projID != stage
}

// IsStage returns true if it's stage env.
func IsStage(projID string) bool {
	return projID == prod
}

// IsProd returns true if it's prod env.
func IsProd(projID string) bool {
	return projID == prod
}
