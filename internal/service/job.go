package service

import (
	"context"
	"fmt"
	"sync"
	"time"
)

type Job struct {
	ID             int64     `json:"id"`
	URL            string    `json:"url"`
	CronExpression string    `json:"cron_expression"`
	Comment        string    `json:"comment"`
	CreatedAt      time.Time `json:"created_at"`
	NextRun        time.Time `json:"next_run"`
	LastRun        time.Time `json:"last_run"`
	LastStatus     int32     `json:"last_status"`
}

type JobService struct {
	jobs      map[int64]*Job
	nextID    int64
	mutex     sync.Mutex
	publisher *EventPublisher
}

func NewJobService(publisher *EventPublisher) *JobService {
	return &JobService{
		jobs:      make(map[int64]*Job),
		nextID:    1,
		publisher: publisher,
	}
}

func (s *JobService) CreateJob(ctx context.Context, url, cronExpression, comment string) (*Job, error) {
	// use go-co-op/gocron + save to kafka
	s.mutex.Lock()
	defer s.mutex.Unlock()

	job := &Job{
		ID:             s.nextID,
		URL:            url,
		CronExpression: cronExpression,
		Comment:        comment,
		CreatedAt:      time.Now(),
		NextRun:        time.Now().Add(5 * time.Minute), // Placeholder
	}

	if s.publisher != nil {
		s.publisher.PublishJobCreated(ctx, job)
	}

	s.jobs[s.nextID] = job
	s.nextID++

	return job, nil
}

func (s *JobService) GetJob(ctx context.Context, id int64) (*Job, bool) {
	job, exists := s.jobs[id]

	return job, exists
}

func (s *JobService) DeleteJob(ctx context.Context, id int64) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	if _, exists := s.jobs[id]; !exists {
		return fmt.Errorf("job with id %d not found", id)
	}

	delete(s.jobs, id)

	if s.publisher != nil {
		s.publisher.PublishJobDeleted(ctx, id)
	}

	return nil
}
