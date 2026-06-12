//HAS THE GO API SERVER
//wire everything here, sliding window and stuff

package main

import (
	"log"
	"net/http"
	"time"
	"os"		//to make port configurable
	"api-gateway/handlers"
	"api-gateway/middleware"
	"api-gateway/redis"		//import redis

	

)

func main(){		//needs a main func always

	redis.ConnectRedis()		//call this function from redis/client.go

	redis.InitCircuitBreaker()

	redis.RedisHealthy.Store(true)
	redis.StartHealthMonitor()

	//create limiter
	limiter := redis.NewRedisSlidingWindowLimiter(
		100,		//100 requests per minute
		time.Minute,
	)

	//create handler
	protectedHandler := http.HandlerFunc(
		handlers.ProtectedHandler,
	)

	//wrap the middleware (algo of rate limiter)
	http.Handle(
		"/api/data",
		middleware.SlidingWindowMiddleware(
			limiter,
			protectedHandler,
		),
	)

	http.HandleFunc(	//register health handler route
		"/health",
		handlers.HealthHandler,
	)

	http.HandleFunc(	//register metrics handler route
		"/metrics",
		handlers.MetricsHandler,
	)

	http.HandleFunc(
		"/dashboard",
		handlers.DashboardHandler,
	)

	port := os.Getenv("PORT")

	if port == "" {
		port = "8080"
	}

	log.Println("Server running on :" + port)

	log.Fatal(
		http.ListenAndServe(		//start server
			":"+port,
			nil,
		),
	)
}




