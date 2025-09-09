package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

func LoadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Println(".env file not found, relying on system environment")
	}
}

func GetDBURL() string {
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return ""
	}
	return dbURL
}
