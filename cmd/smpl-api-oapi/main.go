package main

import (
	"context"
	"net/http"

	"github.com/vlladoff/micro-learn/internal/app"
	"github.com/vlladoff/micro-learn/internal/handler"
	api "github.com/vlladoff/micro-learn/internal/server"
	"github.com/vlladoff/micro-learn/internal/service"
	"go.uber.org/fx"
)

func main() {
	fx.New(
		service.ServiceModule,
		handler.HandlerModule,
		app.AppModule,

		fx.Provide(NewHttpServer),

		fx.Invoke(func(*http.Server) {}),
	).Run()
}

func NewHttpServer(lc fx.Lifecycle, server *app.SmplServer) *http.Server {
	handler := api.Handler(server)

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
