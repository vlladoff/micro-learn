package service

import (
	"context"
	"log"

	"go.uber.org/fx"
)

func StartJobExecutor(lc fx.Lifecycle, executor *JobExecutor) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				if err := executor.Start(context.Background()); err != nil {
					log.Printf("[ERROR] JobExecutor failed to start: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			executor.Stop()
			return nil
		},
	})
}

var ServiceModule = fx.Module("services",
	fx.Provide(
		NewJobService,
		NewJobExecutor,
		NewEventPublisherService,
	),
	fx.Invoke(StartJobExecutor),
)
