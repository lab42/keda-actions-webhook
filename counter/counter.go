package counter

import (
	"context"
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/labstack/gommon/log"
)

type RedisCounter struct {
	rdb        *redis.Client
	counterKey string
}

func NewRedisCounter(address string, password string, db int) (redisCounter RedisCounter, err error) {
	redisCounter.rdb = redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	redisCounter.counterKey = "runner_count"

	return redisCounter, err
}

func (c RedisCounter) TestConnection() bool {
	// Check the connection
	_, err := c.rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Errorf(fmt.Sprintf("Failed to ping Redis server: %v", err))
		return false
	}

	return true
}

func (c RedisCounter) Increment() bool {
	newValue, err := c.rdb.Incr(context.Background(), c.counterKey).Result()
	if err != nil {
		log.Fatalf("Failed to increment Redis counter: %v", err)
		return false
	}

	log.Infof("Incremented Redis runner counter: %s", c.counterKey)
	log.Infof("Queue lenght: %s", newValue)
	return true
}

func (c RedisCounter) Decrement() bool {
	newValue, err := c.rdb.Decr(context.Background(), c.counterKey).Result()
	if err != nil {
		log.Fatalf("Failed to decrement Redis counter: %s", c.counterKey)
		return false
	}

	log.Infof("Decremented Redis runner counter: %s", c.counterKey)
	log.Infof("Queue lenght: %s", newValue)
	return true
}
