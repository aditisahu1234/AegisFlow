package redis //file belongs to redis package

import (
	"api-gateway/config"
	"api-gateway/telemetry"
	"context" //redis operations need a context
	"fmt"
	"log"
	"os"
	"time"

	goredis "github.com/redis/go-redis/v9" //go redis, redist client library
)

var Client *goredis.Client //creates one redis client, later, middleware,
// handlers, rate limiter, all use the same client
var SlidingWindowScript string //script loader

// establishes connection with redis
func ConnectRedis(cfg config.Config) error {

	//no hardcoded address Now your application can run almost anywhere with the same executable.
	//Ask the operating system if an environment variable named REDIS_HOST exists

	redisHost := cfg.RedisHost
	redisPort := cfg.RedisPort
	//build dynamically
	redisAddr := fmt.Sprintf("%s:%s", redisHost, redisPort)

	log.Println("Connecting to Redis at:", redisAddr)

	Client = goredis.NewClient(&goredis.Options{
		Addr: redisAddr,
	})

	//no timeout, or canellation
	ctx, cancel := context.WithTimeout(
		context.Background(),
		3*time.Second,
	)
	defer cancel()

	err := Client.Ping(ctx).Err() //Go--PING--Redis--PONG
	if err != nil {               //if redis responds, err==nil
		RedisHealthy.Store(false)
		telemetry.RedisHealth.Store(0)
		return err //if redis down, err!=nil
	}

	log.Println("Connected to Redis") //load lua once server starts, keep in memory
	RedisHealthy.Store(true)
	telemetry.RedisHealth.Store(1)

	script, erro := os.ReadFile(
		"redis/scripts/sliding_window.lua",
	)

	if erro != nil {
		panic(erro)
	}

	SlidingWindowScript = string(script)
	log.Println("Lua Script Loaded")
	return nil
}
