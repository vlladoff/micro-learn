package config

import "go.uber.org/fx"

var ConfigModule = fx.Module("config",
	fx.Provide(Load),
)