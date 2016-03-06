package follow

import "time"

type Follow struct {
	FollowDateTime time.Time `json:"follow_datetime"`
	FollowerID     string    `json:"follower_id"`
	FollowingID    string    `json:"following_id"`
}
