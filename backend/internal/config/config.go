package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Env            string
	Port           int
	MongoURI       string
	MongoDBName    string
	JWTSecret      string
	JWTIssuer      string
	AccessTokenTTL time.Duration
	DBTimeout      time.Duration
}

func Load() (Config, error) {
	cfg := Config{
		Env:            getEnv("ENV", "development"),
		Port:           getEnvInt("PORT", 8080),
		MongoURI:       getEnv("MONGO_URI", "mongodb://mongo:27017"),
		MongoDBName:    getEnv("MONGO_DB_NAME", "akiba"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTIssuer:      getEnv("JWT_ISSUER", "akiba-api"),
		AccessTokenTTL: getEnvDuration("ACCESS_TOKEN_TTL", time.Hour),
		DBTimeout:      getEnvDuration("DB_TIMEOUT", 5*time.Second),
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET cannot be empty")
	}
	if cfg.Port <= 0 {
		return Config{}, fmt.Errorf("PORT must be > 0")
	}
	return cfg, nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func getEnvInt(k string, def int) int {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return def
	}
	return n
}
func getEnvDuration(k string, def time.Duration) time.Duration {
	v := os.Getenv(k)
	if v == "" {
		return def
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return def
	}
	return d
}
