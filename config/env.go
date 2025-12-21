package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

// LoadEnv loads .env file if present and provides getters
func LoadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using environment variables")
    }
}

// GetEnv returns environment variable value or fallback if not set
func GetEnv(key, fallback string) string {
    if value, ok := os.LookupEnv(key); ok {
        return value
    }
    return fallback
}
