//HAS THE GO API SERVER
//wire everything here, sliding window and stuff

package main

import (
	"api-gateway/config"
	"api-gateway/handlers"
	"api-gateway/internal/runtime"
	"api-gateway/middleware"
	"api-gateway/redis"
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

func main() {
	cfg := config.Load()

	telemetry.InitTelemetry()
	if err := telemetry.InitMetrics(); err != nil {
		log.Fatal(err)
	}
	if err := telemetry.InitTracer(); err != nil {
		log.Fatal(err)
	}

	// Create the runtime manager
	runtimemanager := runtime.NewManager()

	// Register the Redis component (and later other dependencies)
	redisComponent := redis.NewComponent(cfg)
	runtimemanager.Register(redisComponent)

	// Start the health engine in the background
	runtimeCtx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// Start all components (with retries handled inside)
	if err := runtimemanager.Start(runtimeCtx); err != nil {
		log.Fatalf("Failed to start runtime: %v", err)
	}
	runtimemanager.StartHealthEngine(runtimeCtx)

	// Start HTTP server as usual
	limiter := redis.NewRedisSlidingWindowLimiter(100, time.Minute)
	protectedHandler := http.HandlerFunc(handlers.ProtectedHandler)

	http.Handle("/api/data", middleware.SlidingWindowMiddleware(limiter, protectedHandler))

	http.HandleFunc(
		"/health",
		handlers.HealthHandler,
	)

	http.HandleFunc(
		"/ready",
		handlers.ReadinessHandler(runtimemanager),
	)

	http.Handle(
		"/metrics",
		promhttp.Handler(),
	)

	http.HandleFunc(
		"/dashboard",
		handlers.DashboardHandler,
	)

	http.HandleFunc(
		"/live",
		handlers.LivenessHandler,
	)
	server := &http.Server{
		Addr:    ":" + cfg.Port,
		Handler: nil,
	}

	go func() {
		log.Println("Server running on :" + cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server failed: %v", err)
		}
	}()

	<-runtimeCtx.Done()
	log.Println("Shutdown signal received...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Printf("Server shutdown error: %v", err)
	}

	// Gracefully stop all components via runtime manager if you add graceful shutdown
	if err := runtimemanager.Stop(shutdownCtx); err != nil {
		log.Printf("Runtime shutdown error: %v", err)
	}

	log.Println("Server stopped gracefully.")
}
