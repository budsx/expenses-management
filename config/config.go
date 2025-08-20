package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	Database            Database
	Log                 Log
	ServicePort         int
	PaymentProcessorURL string
	JWTKey              string
}

type Database struct {
	DriverName string
	Host       string
	Port       int
	Username   string
	Password   string
	DbName     string
}

type Log struct {
	Level int
}

func Load() *Config {
	godotenv.Load()
	return &Config{
		Database: Database{
			DriverName: getEnv("DB_DRIVER", "postgres"),
			Host:       getEnv("DB_HOST", "localhost"),
			Port:       getEnvInt("DB_PORT", 5432),
			Username:   getEnv("DB_USERNAME", "postgres"),
			Password:   getEnv("DB_PASSWORD", "postgres"),
			DbName:     getEnv("DB_NAME", "postgres"),
		},
		Log: Log{
			Level: getEnvInt("LOG_LEVEL", -1),
		},
		ServicePort:         getEnvInt("SERVICE_PORT", 8000),
		PaymentProcessorURL: getEnv("PAYMENT_PROCESSOR_URL", ""),
		JWTKey:              getEnv("JWT_KEY", ""),
	}
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return defaultValue
	}
	return result
}

func getEnvBool(key string, defaultValue bool) bool {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value == "true"
}
