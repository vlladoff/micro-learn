package main

import (
	"context"
	"net/http"
	"os"
	"strings"

	"github.com/vlladoff/micro-learn/internal/app"
	"github.com/vlladoff/micro-learn/internal/handler"
	"github.com/vlladoff/micro-learn/internal/infrastructure/kafka"
	"github.com/vlladoff/micro-learn/internal/middleware"
	api "github.com/vlladoff/micro-learn/internal/server"
	"github.com/vlladoff/micro-learn/internal/service"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		fx.Provide(
			NewKafkaProducer,
			NewKafkaConsumer,
		),

		service.ServiceModule,
		handler.HandlerModule,
		middleware.MiddlewareModule,
		app.AppModule,

		fx.Provide(NewHttpServer),

		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func NewHttpServer(lc fx.Lifecycle, server *app.SmplServer) *http.Server {
	handler := api.HandlerWithOptions(server,
		api.ChiServerOptions{
			Middlewares: []api.MiddlewareFunc{
				middleware.AddRequestId,
			},
		})

	srv := &http.Server{
		Addr:    ":8080",
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

func NewKafkaProducer() (*kafka.Producer, error) {
	brokers := getBrokers()
	return kafka.NewProducer(brokers)
}

func NewKafkaConsumer() (*kafka.Consumer, error) {
	brokers := getBrokers()
	groupID := getEnv("KAFKA_GROUP_ID", "job-service")
	return kafka.NewConsumer(brokers, groupID)
}

func getBrokers() []string {
	brokers := getEnv("KAFKA_BROKERS", "localhost:9092")
	return strings.Split(brokers, ",")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
