package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/vlladoff/micro-learn/internal/app"
	"github.com/vlladoff/micro-learn/internal/config"
	"github.com/vlladoff/micro-learn/internal/handler"
	"github.com/vlladoff/micro-learn/internal/infrastructure/kafka"
	"github.com/vlladoff/micro-learn/internal/infrastructure/redis"
	"github.com/vlladoff/micro-learn/internal/middleware"
	"github.com/vlladoff/micro-learn/internal/repository"
	api "github.com/vlladoff/micro-learn/internal/server"
	"github.com/vlladoff/micro-learn/internal/service"
	"go.uber.org/fx"
)

func main() {
	if len(os.Args) > 1 && os.Args[1] == "healthcheck" {
		healthCheck()
		return
	}

	fx.New(
		config.ConfigModule,
		redis.RedisModule,
		kafka.KafkaModule,
		repository.RepositoryModule,
		service.ServiceModule,
		handler.HandlerModule,
		middleware.MiddlewareModule,
		app.AppModule,

		fx.Provide(NewHttpServer),

		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func NewHttpServer(lc fx.Lifecycle, server *app.SmplServer, cfg *config.Config) *http.Server {
	handler := api.HandlerWithOptions(server, api.ChiServerOptions{
		Middlewares: []api.MiddlewareFunc{
			middleware.RequestIDMiddleware,
			middleware.AuthMiddleware,
		},
	})

	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	srv := &http.Server{
		Addr:    addr,
		Handler: handler,
	}

	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go srv.ListenAndServe()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			return srv.Shutdown(ctx)
		},
	})

	return srv
}

func healthCheck() {
	cfg, err := config.Load()
	if err != nil {
		log.Printf("Failed to load config for healthcheck: %v", err)
		os.Exit(1)
	}

	client := &http.Client{Timeout: 3 * time.Second}
	url := fmt.Sprintf("http://%s:%d/ping", cfg.Server.Host, cfg.Server.Port)

	resp, err := client.Get(url)
	if err != nil {
		log.Printf("Health check failed: %v", err)
		os.Exit(1)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Health check failed: status %d", resp.StatusCode)
		os.Exit(1)
	}

	log.Println("Health check passed")
}
