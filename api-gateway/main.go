//HAS THE GO API SERVER
//wire everything here, sliding window and stuff

package main

import (
	"api-gateway/config"
	"api-gateway/handlers"
	"api-gateway/middleware"
	"api-gateway/redis" //import redis
	"api-gateway/telemetry"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() { //needs a main func always

	cfg := config.Load()

	telemetry.InitTelemetry()
	if err := telemetry.InitMetrics(); err != nil {
		log.Fatal(err)
	}
	if err := telemetry.InitTracer(); err != nil {
		log.Fatal(err)
	}

	//call this function from redis/client.go
	redis.ConnectRedis(cfg)

	redis.InitCircuitBreaker()

	redis.RedisHealthy.Store(true)
	telemetry.RedisHealth.Store(1)

	redis.StartHealthMonitor()

	//create limiter
	limiter := redis.NewRedisSlidingWindowLimiter(
		100, //100 requests per minute
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

	http.HandleFunc( //register health handler route
		"/health",
		handlers.HealthHandler,
	)
	/*
		http.HandleFunc( //register metrics handler route
			"/metrics-json",
			handlers.MetricsHandler,
		)
	*/
	http.Handle( //adding new route
		"/metrics",
		promhttp.Handler(),
	)

	http.HandleFunc(
		"/dashboard",
		handlers.DashboardHandler,
	)

	log.Println("Server running on :" + cfg.Port)

	log.Fatal(
		http.ListenAndServe( //start server
			":"+cfg.Port,
			nil,
		),
	)
}
