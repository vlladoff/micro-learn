package service

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Job struct {
	ID             string    `json:"id"`
	URL            string    `json:"url"`
	CronExpression string    `json:"cron_expression"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"created_at"`
	NextRun        time.Time `json:"next_run"`
	LastRun        time.Time `json:"last_run,omitempty"`
	LastStatus     int32     `json:"last_status,omitempty"`
}

type JobService struct {
	jobs      map[string]*Job
	mutex     sync.Mutex
	publisher *EventPublisher
}

func NewJobService(publisher *EventPublisher) *JobService {
	return &JobService{
		jobs:      make(map[string]*Job),
		publisher: publisher,
	}
}

func (s *JobService) CreateJob(ctx context.Context, url, cronExpression, comment string) (*Job, error) {
	// use go-co-op/gocron + save to kafka
	s.mutex.Lock()
	defer s.mutex.Unlock()

	jobID := uuid.New().String()

	job := &Job{
		ID:             jobID,
		URL:            url,
		CronExpression: cronExpression,
		Comment:        comment,
		CreatedAt:      time.Now(),
		// NextRun:        parse with go-co-op/gocron
	}

	s.jobs[jobID] = job

	if s.publisher != nil {
		s.publisher.PublishJobCreated(ctx, job)
	}

	return job, nil
}

func (s *JobService) GetJob(ctx context.Context, id string) (*Job, bool) {
	job, exists := s.jobs[id]

	return job, exists
}

func (s *JobService) DeleteJob(ctx context.Context, id string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.jobs[id]; !exists {
		return fmt.Errorf("job with id %s not found", id)
	}

	delete(s.jobs, id)

	if s.publisher != nil {
		s.publisher.PublishJobDeleted(ctx, id)
	}

	return nil
}
