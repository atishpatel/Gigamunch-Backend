package session

import (
	"encoding/json"
	"testing"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"

	"google.golang.org/appengine/aetest"
)

func TestSaveUserSession(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()

	var nilUser *types.User

	testCases := []struct {
		description string
		uuid        string
		user        *types.User
		output      error
	}{
		{
			description: "Invalid UUID",
			uuid:        "klj",
			user:        &types.User{},
			output:      errors.ErrInvalidUUID,
		},
		{
			description: "nil for user",
			uuid:        "b4e4f890-2210-4ff3-a67b-60be9989ce68",
			user:        nilUser,
			output:      errors.ErrNilParamenter,
		},
		{
			description: "Save a user",
			uuid:        "b4e4f890-2210-4ff3-a67b-60be9989ce68",
			user: &types.User{
				Email:       "test@test.com",
				Name:        "name",
				PhotoURL:    "url",
				Permissions: 0,
			},
			output: nil,
		},
	}
	config := config.GetConfig()
	redisClient := getRedisClient(config.RedisSessionServerIP, config.RedisSessionServerPassword, 0, 1)
	_, err = redisClient.Ping().Result()
	if err != nil {
		t.Skip("Failed to connect to redis client. Skipping the SaveUserSession test.")
	}

	for _, test := range testCases {
		errChan := SaveUserSession(ctx, test.uuid, test.user)
		err = <-errChan
		if err != test.output {
			t.Errorf("Failed test %s | expected error: %+v | got error: %+v", test.description, test.output, err)
		}
		if err == nil {
			serialized, err := redisClient.Get(SessionNamespace + test.uuid).Result()
			if err != nil {
				t.Errorf("Failed to get from redis %+v", err)
			}
			testUser := &types.User{}
			err = json.Unmarshal([]byte(serialized), testUser)
			if err != nil {
				t.Errorf("Failed to unmarshal data %+v", err)
			}
			if testUser.Email != test.user.Email || testUser.Name != test.user.Name ||
				testUser.PhotoURL != test.user.PhotoURL ||
				testUser.Permissions != test.user.Permissions {
				t.Errorf("Failed test %s | saved user and expected user do not match", test.description)
			}
		}
	}
}
