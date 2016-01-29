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

	// setup
	var nilUser *types.User
	config := config.GetConfig()
	redisClient := getRedisClient(config.RedisSessionServerIP, config.RedisSessionServerPassword, 0, 1)
	_, err = redisClient.Ping().Result()
	if err != nil {
		t.Skip("Failed to connect to redis client")
	}
	testCases := []struct {
		description string
		uuid        string
		user        *types.User
		output      error
	}{
		{
			description: "Invalid UUID",
			uuid:        "klj",
			user:        &types.User{}, // sets all variables to defaults
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
				Email: "test@test.com",
			},
			output: nil,
		},
	}
	// run test
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
			// Failed some test
			if *testUser != *test.user {
				t.Errorf("Failed test %s | saved user and expected user do not match", test.description)
			}
		}
	}
}

func TestGetUserSession(t *testing.T) {
	ctx, done, err := aetest.NewContext()
	if err != nil {
		t.Fatal(err)
	}
	defer done()
	// setup
	config := config.GetConfig()
	redisClient := getRedisClient(config.RedisSessionServerIP, config.RedisSessionServerPassword, 0, 1)
	_, err = redisClient.Ping().Result()
	if err != nil {
		t.Skip("Failed to connect to redis client")
	}

	validUUID := "b4e4f890-2210-4ff3-a67b-60be9989ce68"
	expectedValidUser := &types.User{
		Email:       "test@test.com",
		Name:        "name",
		PhotoURL:    "url",
		Permissions: 0,
	}
	SaveUserSession(ctx, validUUID, expectedValidUser)
	testCases := []struct {
		description string
		uuid        string
		output      *types.User
	}{
		{
			description: "Invalid UUID",
			uuid:        "klj",
			output:      nil,
		},
		{
			description: "Getting a user that doesn't exist",
			uuid:        "b4e4f890-2210-4ff3-a67b-60be9989ce67",
			output:      nil,
		},
		{
			description: "Getting a valid user",
			uuid:        validUUID,
			output:      expectedValidUser,
		},
	}
	// run test
	for _, test := range testCases {
		userChan := GetUserSession(ctx, test.uuid)
		user := <-userChan
		if user != nil && test.output != nil {
			if *user != *test.output {
				t.Errorf("Failed test %s | expected error: %+v | got error: %+v", test.description, test.output, err)
			}
		}
	}
}
