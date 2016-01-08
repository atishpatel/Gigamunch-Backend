package databases

import "gopkg.in/redis.v3"

type Config struct {
	RedisSessionDBIP       string
	RedisSessionDBPassword string
}

type Database struct {
	redisSessionClient *redis.Client
}

func CreateDatabase(config Config) *Database {
	db := Database{}
	db.redisSessionClient = createRedisDatabase(config.RedisSessionDBIP, config.RedisSessionDBPassword)
	return &db
}
