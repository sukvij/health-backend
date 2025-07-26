package redisservice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
)

func NewRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: "localhost:6379",
		// docker run
		// Addr:     "redis:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	_, err := client.Ping().Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}
	return client
}

func SetValueToRedis(redisClient *redis.Client, key string, value interface{}) error {

	jsonValue, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %v", err)
	}
	err = redisClient.Set(key, jsonValue, 3600*time.Second).Err()
	if err != nil {
		return fmt.Errorf("set key problem for val %v and err is %v", value, err)
	}
	return redis.Nil
}

func GetValueFromRedis(redisClient *redis.Client, key string) (interface{}, error) {
	byteVal, err := redisClient.Get(key).Result()
	if err != nil {
		return "", fmt.Errorf("get key problem for key %v and err is %v", key, err)
	}
	var res interface{}
	json.Unmarshal([]byte(byteVal), &res)
	fmt.Println("val is", res)
	return res, redis.Nil
}
