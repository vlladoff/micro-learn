package service

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
	"github.com/vlladoff/micro-learn/internal/repository"
)

type JobService struct {
	repository repository.JobRepository
	publisher  *EventPublisher
}

func NewJobService(publisher *EventPublisher, jobRepository repository.JobRepository) *JobService {
	return &JobService{
		repository: jobRepository,
		publisher:  publisher,
	}
}

func (s *JobService) CreateJob(ctx context.Context, url, cronExpression, comment string) (*repository.Job, error) {
	jobID := uuid.New().String()

	job := &repository.Job{
		ID:             jobID,
		URL:            url,
		CronExpression: cronExpression,
		Comment:        comment,
		CreatedAt:      time.Now(),
	}

	if err := s.repository.Save(ctx, job); err != nil {
		return nil, fmt.Errorf("failed to save job: %w", err)
	}

	if s.publisher != nil {
		if err := s.publisher.PublishJobCreated(ctx, job); err != nil {
			log.Printf("[ERROR] Failed to publish job created event: %v", err)
		}
	}

	return job, nil
}

func (s *JobService) GetJob(ctx context.Context, id string) (*repository.Job, bool) {
	job, err := s.repository.FindByID(ctx, id)
	if err != nil {
		log.Printf("[ERROR] Failed to get job: %v", err)
		return nil, false
	}

	return job, true
}

func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	if _, exists := s.GetJob(ctx, id); !exists {
		return fmt.Errorf("job with id %s not found", id)
	}

	if err := s.repository.Delete(ctx, id); err != nil {
		return fmt.Errorf("failed to delete job: %w", err)
	}

	if s.publisher != nil {
		if err := s.publisher.PublishJobDeleted(ctx, id); err != nil {
			log.Printf("[ERROR] Failed to publish job deleted event: %v", err)
		}
	}

	return nil
}
