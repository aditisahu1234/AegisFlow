//HAS THE GO API SERVER
//wire everything here, sliding window and stuff

package main

import (
	"api-gateway/config"
	"api-gateway/handlers"
	"api-gateway/middleware"
	"api-gateway/redis" //import redis
	"api-gateway/telemetry"

	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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
	startupCtx := context.Background()
	go redis.ConnectWithRetry(
		startupCtx,
		cfg,
	)

	redis.InitCircuitBreaker()

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

	// Create HTTP server
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: nil,
	}

	// Start server in a separate goroutine
	go func() {
		log.Println("Server running on :" + cfg.Port)

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	// Wait for shutdown signal
	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)

	defer stop()

	<-ctx.Done()

	log.Println("Shutdown signal received...")

	// Allow up to 10 seconds for graceful shutdown
	shutdownCtx, cancel := context.WithTimeout(
		context.Background(),
		10*time.Second,
	)
	defer cancel()

	// Gracefully stop HTTP server
	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Close Redis connection
	if redis.Client != nil {
		redis.Client.Close()
	}

	log.Println("Server stopped gracefully.")
}
