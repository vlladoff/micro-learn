package handler

import "go.uber.org/fx"

var HandlerModule = fx.Module("handlers",
	fx.Provide(
		NewDefaultHandler,
		NewJobHandler,
	),
)
