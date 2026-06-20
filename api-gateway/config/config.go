//will read and contian all environment variables

package config

import "os"

type Config struct {
	Port          string
	RedisHost     string
	RedisPort     string
	RedisPassword string

	Environment string

	RequestTimeout string
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func Load() Config {
	return Config{
		Port:      getEnv("PORT", "8080"),
		RedisHost: getEnv("REDIS_HOST", "localhost"),
		RedisPort: getEnv("REDIS_PORT", "6379"),
	}
}
