package redis	//file belongs to redis package

import (
	"context"			//redis operations need a context
	"fmt"
	"os"
	goredis "github.com/redis/go-redis/v9"		//go redis, redist client library
)

var Client *goredis.Client		//creates one redis client, later, middleware,
								//handlers, rate limiter, all use the same client
var SlidingWindowScript string		//script loader


//establishes connection with redis
func ConnectRedis() {
	Client = goredis.NewClient(&goredis.Options{	//connects go and redis client to local host
		Addr: "localhost:6379",
	})				
	
	ctx:=context.Background() //no timeout, or canellation

	err:=Client.Ping(ctx).Err()		//Go--PING--Redis--PONG
	if err != nil {		//if redis responds, err==nil
		panic(err)		//if redis down, err!=nil
	}

	fmt.Println("Connected to Redis")		//load lua once server starts, keep in memory
	script, erro := os.ReadFile(
		"redis/scripts/sliding_window.lua",
	)
	
	if erro != nil {		
		panic(erro)		
	}
	
	SlidingWindowScript = string(script)
	fmt.Println("Lua Script Loaded")
}