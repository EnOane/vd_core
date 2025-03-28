package config

import (
	"github.com/joho/godotenv"
	"github.com/rs/zerolog/log"
	"os"
	"strconv"
)

type natsConfig struct {
	Host string
	Port int
}

type Config struct {
	NatsConfig natsConfig
}

var AppConfig Config

func MustLoad() {
	if err := godotenv.Load(); err != nil {
		log.Fatal().Msg("No .env file found")
	}

	AppConfig = Config{
		NatsConfig: natsConfig{
			Host: getEnv("NATS_HOST", ""),
			Port: getEnvAsInt("NATS_PORT", 4222),
		},
	}
}

func getEnv(key string, defaultVal string) string {
	if v, ok := os.LookupEnv(key); ok {
		return v
	}

	return defaultVal
}

func getEnvAsInt(key string, defaultVal int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}

	return defaultVal
}
