package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	App_Env         string
	Port            string
	DatabaseURL     string
	JWTSecret       string
	RedisAddr       string
	RedisPassword   string
	RedisDB         int
	Environment     string
	StorageProvider string

	
    // // S3 (only needed if STORAGE_PROVIDER=s3)
    // S3Bucket string
    // S3Region string
    // S3AccessKey string
    // S3SecretKey string
}

var Config *AppConfig

func LoadEnv() {
	_ = godotenv.Load()

	dbStr := getEnv("REDIS_DB", "0")
	db, _ := strconv.Atoi(dbStr)

	Config = &AppConfig{
		App_Env:         getEnv("APP_ENV", "local"),
		Port:            getEnv("PORT", "8080"),
		DatabaseURL:     getEnv("DATABASE_PATH", ""),
		JWTSecret:       getEnv("JWT_SECRET", "yourjwtsecret"),
		RedisAddr:       getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:   os.Getenv("REDIS_PASSWORD"),
		RedisDB:         db,
		StorageProvider: getEnv("STORAGE_PROVIDER", "local"),
	}

	log.Printf(" [env=%s]", Config.Environment)
}

func getEnv(key, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}
