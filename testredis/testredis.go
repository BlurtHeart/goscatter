package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

type RedisClient struct {
	Addr     string
	Password string
	DB       int
	Client   *redis.Client
}

func (rc *RedisClient) Connect() {
	if rc.Client != nil {
		pong, err := rc.Client.Ping().Result()
		if err == nil && pong == "PONG" {
			return
		}
	}
	rc.Client = redis.NewClient(&redis.Options{
		Addr:     rc.Addr,
		Password: rc.Password,
		DB:       rc.DB,
	})
}

func (rc *RedisClient) SetString(key, val string) error {
	err := rc.Client.Set(key, val, 0).Err()
	return err
}

func (rc *RedisClient) GetString(key string) (string, error) {
	return rc.Client.Get(key).Result()
}

func (rc *RedisClient) SetHash(key string, fields map[string]interface{}) error {
	return rc.Client.HMSet(key, fields).Err()
}

func (rc *RedisClient) GetHash(key string) (map[string]string, error) {
	return rc.Client.HGetAll(key).Result()
}

func (rc *RedisClient) SetList(key, val string) error {
	return rc.Client.LPush(key, val).Err()
}

func (rc *RedisClient) GetList(key string, start, stop int64) ([]string, error) {
	return rc.Client.LRange(key, start, stop).Result()
}

func main() {
	client := &RedisClient{
		Addr:     "127.0.0.1:6379",
		Password: "",
		DB:       0,
	}
	client.Connect()

	// value is string
	key := "user11"
	val := "val11"
	err := client.SetString(key, val)
	if err != nil {
		panic(err)
	}
	val2, err := client.GetString(key)
	if err != nil {
		panic(err)
	}
	fmt.Println("key:", key, "val:", val2)

	// value is hash
	key2 := "user12"
	fields := make(map[string]interface{})
	fields["username"] = "abc"
	fields["age"] = 15
	fields["score"] = 97.5
	err = client.SetHash(key2, fields)
	if err != nil {
		panic(err)
	}
	fields2, err := client.GetHash(key2)
	if err != nil {
		panic(err)
	}
	fmt.Println(fields2)

	// value is list
	key3 := "user13"
	valuelist := []string{"test1", "test2", "test3"}
	for _, v := range valuelist {
		client.SetList(key3, v)
	}
	valuelist2, err := client.GetList(key3, 0, -1)
	if err != nil {
		panic(err)
	}
	fmt.Println(valuelist2)

	// fuzzy query
	keys, err := client.Client.Keys("test*").Result()
	if err != nil {
		panic(err)
	}
	for _, k := range keys {
		v, err := client.GetString(k)
		if err != nil {
			panic(err)
		}
		fmt.Println("key:", k, "value:", v)
	}

}
