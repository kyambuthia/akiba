package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
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

	db := client.Database(cfg.MongoDBName)
	userRepo := mongoRepo.NewUserRepository(db, cfg.DBTimeout)
	indexCtx, indexCancel := context.WithTimeout(context.Background(), cfg.DBTimeout)
	defer indexCancel()
	if err := userRepo.EnsureIndexes(indexCtx); err != nil {
		log.Fatalf("index setup error: %v", err)
	}

	jwtMgr := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer)
	authSvc := usecase.NewAuthService(userRepo, jwtMgr, cfg.AccessTokenTTL)
	router := httptransport.NewRouter(logger, authSvc, jwtMgr, func(ctx context.Context) error {
		return client.Ping(ctx, nil)
	})

	srv := &http.Server{Addr: fmt.Sprintf(":%d", cfg.Port), Handler: router, ReadHeaderTimeout: 5 * time.Second}
	logger.Info("starting api", "env", cfg.Env, "port", cfg.Port)

	serverErr := make(chan error, 1)
	go func() {
		serverErr <- srv.ListenAndServe()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	select {
	case sig := <-stop:
		logger.Info("shutdown signal received", "signal", sig.String())
	case err := <-serverErr:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("server error: %v", err)
		}
		return
	}

	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("server shutdown failed", "error", err)
	}
	if err := client.Disconnect(shutdownCtx); err != nil {
		logger.Error("mongo disconnect failed", "error", err)
	}
}
