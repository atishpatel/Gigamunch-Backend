package session

import (
	"encoding/json"
	"time"

	"golang.org/x/net/context"

	"gopkg.in/redis.v3"

	"log"

	"github.com/atishpatel/Gigamunch-Backend/config"
	"github.com/atishpatel/Gigamunch-Backend/errors"
	"github.com/atishpatel/Gigamunch-Backend/types"
	"github.com/atishpatel/Gigamunch-Backend/utils"
)

const (
	// SessionNamespace is pre-fixed to a session
	SessionNamespace = "session:"
	// UserSessionListNamespace is pre-fixed to a list of session based on email
	UserSessionListNamespace = SessionNamespace + ":user:"
)

var (
	redisSessionClient *redis.Client
)

// SaveUserSession saves a types.User in the session db
// if uuid is invalid, errors.ErrInvalidUUID is returned
// if user is nil, errors.ErrNilParamenter is returned
func SaveUserSession(ctx context.Context, UUID string, user *types.User) <-chan error {
	errChannel := make(chan error)
	go func(ctx context.Context, UUID string, user *types.User) {
		defer close(errChannel)
		if !utils.IsValidUUID(UUID) {
			errChannel <- errors.ErrInvalidUUID
			return
		}
		if user == nil {
			errChannel <- errors.ErrNilParamenter
			return
		}
		// log time for request
		startTime := time.Now()
		defer func() { utils.Latencyf(ctx, utils.REDISLATENCY, "SaveUserSession %d", time.Since(startTime)) }()
		// serialize data
		serialized, err := json.Marshal(&user)
		if err != nil {
			errChannel <- err
			return
		}
		// TODO(Atish): save a list of sessionID queryable by email
		err = redisSessionClient.Set(SessionNamespace+UUID, serialized, 0).Err()
		errChannel <- err
	}(ctx, UUID, user)
	return errChannel
}

// GetUserSession gets a types.User from the session db
// if uuid is invalid or user does not exist, nil is returned
func GetUserSession(ctx context.Context, UUID string) <-chan *types.User {
	userChannel := make(chan *types.User)
	go func(ctx context.Context, UUID string) {
		defer close(userChannel)
		if !utils.IsValidUUID(UUID) {
			userChannel <- nil
			return
		}
		// log time for request
		startTime := time.Now()
		defer func() { utils.Latencyf(ctx, utils.REDISLATENCY, "GetUserSession %d", time.Since(startTime)) }()
		// get the serialized user
		serialized, err := redisSessionClient.Get(SessionNamespace + UUID).Result()
		if err != nil {
			utils.Errorf(ctx, "Error getting user from session db: %+v", err)
			userChannel <- nil
			return
		}
		// deserialize user
		user := &types.User{}
		err = json.Unmarshal([]byte(serialized), user)
		if err != nil {
			utils.Errorf(ctx, "Error unmarshaling user: %+v", err)
			userChannel <- nil
			return
		}
		userChannel <- user
	}(ctx, UUID)
	return userChannel
}

func getRedisClient(address string, password string, db int64, poolsize int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
		PoolSize: poolsize,
	})
}

func init() {
	config := config.GetConfig()
	redisSessionClient = getRedisClient(config.RedisSessionServerIP, config.RedisSessionServerPassword, 0, 10)
	_, err := redisSessionClient.Ping().Result()
	if err != nil {
		log.Fatal(err)
	}
}
