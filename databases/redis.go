package databases

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"

	"gopkg.in/redis.v3"

	"log"

	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

const (
	SESSIONNAMESPACE = "session"
)

// SaveUserSession saves a types.User in the session db
func (db *Database) SaveUserSession(ctx context.Context, UUID string, user *types.User) <-chan error {
	errChannel := make(chan error)
	go func(ctx context.Context, UUID string, user *types.User) {
		// log time for request
		startTime := time.Now()
		defer func() { utils.Latencyf(ctx, utils.REDISLATENCY, "SaveUserSession %d", time.Since(startTime)) }()
		defer close(errChannel)
		// serialize data
		serialized, err := json.Marshal(&user)
		if err != nil {
			errChannel <- err
		}
		err = db.redisSessionClient.Set(SESSIONNAMESPACE+UUID, serialized, 0).Err()
		errChannel <- err
	}(ctx, UUID, user)
	return errChannel
}

// GetUserSession gets a types.User from the session db
func (db *Database) GetUserSession(ctx context.Context, UUID string) <-chan *types.User {
	userChannel := make(chan *types.User)
	go func(ctx context.Context, UUID string) {
		// log time for request
		startTime := time.Now()
		defer func() { utils.Latencyf(ctx, utils.REDISLATENCY, "GetUserSession %d", time.Since(startTime)) }()
		defer close(userChannel)

		serialized, err := db.redisSessionClient.Get(UUID).Result()
		if err != nil {
			utils.Errorf(ctx, "Error getting user from session db: %+v", err)
			userChannel <- nil
		}
		user := &types.User{}
		err = json.Unmarshal([]byte(serialized), user)
		if err != nil {
			utils.Errorf(ctx, "Error unmarshaling user: %+v", err)
			userChannel <- nil
		}
	}(ctx, UUID)
	return userChannel
}

func createRedisDatabase(ip string, password string) *redis.Client {
	RedisSessionClient := redis.NewClient(&redis.Options{
		Addr:     ip,
		Password: password,
		DB:       0,
	})

	_, err := RedisSessionClient.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}

	return RedisSessionClient
}
