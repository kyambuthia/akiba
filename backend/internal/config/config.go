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
	port, err := getEnvInt("PORT", 8080)
	if err != nil {
		return Config{}, err
	}
	accessTokenTTL, err := getEnvDuration("ACCESS_TOKEN_TTL", time.Hour)
	if err != nil {
		return Config{}, err
	}
	dbTimeout, err := getEnvDuration("DB_TIMEOUT", 5*time.Second)
	if err != nil {
		return Config{}, err
	}

	cfg := Config{
		Env:            getEnv("ENV", "development"),
		Port:           port,
		MongoURI:       getEnv("MONGO_URI", "mongodb://mongo:27017"),
		MongoDBName:    getEnv("MONGO_DB_NAME", "akiba"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTIssuer:      getEnv("JWT_ISSUER", "akiba-api"),
		AccessTokenTTL: accessTokenTTL,
		DBTimeout:      dbTimeout,
	}
	if cfg.JWTSecret == "" {
		return Config{}, fmt.Errorf("JWT_SECRET cannot be empty")
	}
	if cfg.Port <= 0 {
		return Config{}, fmt.Errorf("PORT must be > 0")
	}
	if cfg.AccessTokenTTL <= 0 {
		return Config{}, fmt.Errorf("ACCESS_TOKEN_TTL must be > 0")
	}
	if cfg.DBTimeout <= 0 {
		return Config{}, fmt.Errorf("DB_TIMEOUT must be > 0")
	}
	return cfg, nil
}

func getEnv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}
func getEnvInt(k string, def int) (int, error) {
	v := os.Getenv(k)
	if v == "" {
		return def, nil
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return 0, fmt.Errorf("%s must be an integer: %w", k, err)
	}
	return n, nil
}
func getEnvDuration(k string, def time.Duration) (time.Duration, error) {
	v := os.Getenv(k)
	if v == "" {
		return def, nil
	}
	d, err := time.ParseDuration(v)
	if err != nil {
		return 0, fmt.Errorf("%s must be a valid duration: %w", k, err)
	}
	return d, nil
}
