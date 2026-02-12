package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"akiba/backend/internal/auth"
	"akiba/backend/internal/config"
	mongoRepo "akiba/backend/internal/infrastructure/mongo"
	"akiba/backend/internal/observability"
	httptransport "akiba/backend/internal/transport/http"
	"akiba/backend/internal/usecase"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}
	logger := observability.NewLogger()
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(cfg.MongoURI))
	if err != nil {
		log.Fatalf("mongo connect error: %v", err)
	}
	defer func() { _ = client.Disconnect(context.Background()) }()

	db := client.Database(cfg.MongoDBName)
	userRepo := mongoRepo.NewUserRepository(db, cfg.DBTimeout)
	if err := userRepo.EnsureIndexes(context.Background()); err != nil {
		log.Fatalf("index setup error: %v", err)
	}

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer)
	authSvc := usecase.NewAuthService(userRepo, jwtMgr, cfg.AccessTokenTTL)
	router := httptransport.NewRouter(logger, authSvc, jwtMgr)

	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: router, ReadHeaderTimeout: 5 * time.Second}
	logger.Info("starting api", "env", cfg.Env, "port", cfg.Port)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatalf("server error: %v", err)
	}
}
