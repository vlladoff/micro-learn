package service

import (
	"context"
	"time"

	"github.com/vlladoff/micro-learn/internal/events"
	"github.com/vlladoff/micro-learn/internal/infrastructure/kafka"
)

type EventPublisher struct {
	producer *kafka.Producer
	topic    string
}

func NewEventPublisherService(producer *kafka.Producer) *EventPublisher {
	return &EventPublisher{
		producer: producer,
		topic:    "job-events",
	}
}

func (ep *EventPublisher) PublishJobCreated(ctx context.Context, job *Job) error {
	event := &events.JobEvent{
		Type:      events.JobCreated,
		JobID:     job.ID,
		Timestamp: time.Now(),
		Data: events.JobCreatedData{
			URL:            job.URL,
			CronExpression: job.CronExpression,
			Comment:        job.Comment,
		},
	}

	return ep.producer.PublishEvent(ctx, ep.topic, event)
}

func (ep *EventPublisher) PublishJobDeleted(ctx context.Context, jobID int64) error {
	event := &events.JobEvent{
		Type:      events.JobDeleted,
		JobID:     jobID,
		Timestamp: time.Now(),
	}

	return ep.producer.PublishEvent(ctx, ep.topic, event)
}
