package testutil

import (
	"fmt"
	"os"
	"testing"

	"github.com/go-redis/redis/v8"
)

func OpenRedisForTest(t *testing.T) *redis.Client {
	t.Helper()

	host := "127.0.0.1"
	port := 16379
	if _, defined := os.LookupEnv("CI"); defined {
		port = 6379
	}

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", host, port),
		Password: "",
		DB:       0,
	})
	if err := client.Ping(client.Context()).Err(); err != nil {
		t.Fatalf("failed to ping to redis: %v", err)
	}

	return client
}
