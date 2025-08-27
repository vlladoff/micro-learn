package kafka

import (
	"context"

	"github.com/ThreeDotsLabs/watermill"
	"github.com/ThreeDotsLabs/watermill-kafka/v2/pkg/kafka"
	"github.com/ThreeDotsLabs/watermill/message"
	"github.com/vlladoff/micro-learn/internal/events"
)

type Producer struct {
	publisher message.Publisher
	logger    watermill.LoggerAdapter
}

func NewProducer(brokers []string) (*Producer, error) {
	logger := watermill.NewStdLogger(false, false)

	config := kafka.PublisherConfig{
		Brokers:   brokers,
		Marshaler: kafka.DefaultMarshaler{},
	}

	publisher, err := kafka.NewPublisher(config, logger)
	if err != nil {
		return nil, err
	}

	return &Producer{
		publisher: publisher,
		logger:    logger,
	}, nil
}

func (p *Producer) PublishEvent(ctx context.Context, topic string, event *events.JobEvent) error {
	payload, err := event.ToJSON()
	if err != nil {
		return err
	}

	msg := message.NewMessage(watermill.NewUUID(), payload)
	msg.Metadata.Set("event_type", string(event.Type))
	msg.Metadata.Set("job_id", event.JobID)

	return p.publisher.Publish(topic, msg)
}

func (p *Producer) Close() error {
	return p.publisher.Close()
}

type Consumer struct {
	subscriber message.Subscriber
	logger     watermill.LoggerAdapter
}

func NewConsumer(brokers []string, groupID string) (*Consumer, error) {
	logger := watermill.NewStdLogger(false, false)

	config := kafka.SubscriberConfig{
		Brokers:       brokers,
		Unmarshaler:   kafka.DefaultMarshaler{},
		ConsumerGroup: groupID,
	}

	subscriber, err := kafka.NewSubscriber(config, logger)
	if err != nil {
		return nil, err
	}

	return &Consumer{
		subscriber: subscriber,
		logger:     logger,
	}, nil
}

func (c *Consumer) Consume(ctx context.Context, topic string, handler func(event *events.JobEvent) error) error {
	messages, err := c.subscriber.Subscribe(ctx, topic)
	if err != nil {
		return err
	}

	go func() {
		for msg := range messages {
			var event events.JobEvent
			if err := event.FromJSON(msg.Payload); err != nil {
				c.logger.Error("failed to unmarshal event", err, nil)
				msg.Nack()
				continue
			}

			if err := handler(&event); err != nil {
				c.logger.Error("Failed to handle event", err, nil)
				msg.Nack()
				continue
			}

			msg.Ack()
		}
	}()

	return nil
}

func (c *Consumer) Close() error {
	return c.subscriber.Close()
}
