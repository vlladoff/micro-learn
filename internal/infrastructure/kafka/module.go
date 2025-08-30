package kafka

import "go.uber.org/fx"

var KafkaModule = fx.Module("kafka",
	fx.Provide(
		NewKafkaProducer,
		NewKafkaConsumer,
	),
)
