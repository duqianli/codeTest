package test

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"testing"
	"time"
)

var ct = context.Background()
var drb = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "",
	DB:       0,
})

func TestSetRedis(t *testing.T) {
	drb.Set(ct, "name", "mmc", time.Second*10)
}

func TestGetRedis(t *testing.T) {
	va, err := drb.Get(ct, "name").Result()
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(va)
}
