package service

import "go.uber.org/fx"

var ServiceModule = fx.Module("services",
	fx.Provide(
		NewDefaultService,
		NewJobService,
		NewEventPublisherService,
	),
)
