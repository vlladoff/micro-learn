package service

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-co-op/gocron"
	"github.com/vlladoff/micro-learn/internal/config"
	"github.com/vlladoff/micro-learn/internal/events"
	"github.com/vlladoff/micro-learn/internal/infrastructure/kafka"
	"github.com/vlladoff/micro-learn/internal/repository"
)

type JobExecutor struct {
	repository repository.JobRepository
	scheduler  *gocron.Scheduler
	client     *http.Client
	consumer   *kafka.Consumer
}

func NewJobExecutor(jobRepository repository.JobRepository, kafkaConsumer *kafka.Consumer, cfg *config.Config) *JobExecutor {
	scheduler := gocron.NewScheduler(time.UTC)
	scheduler.StartAsync()

	return &JobExecutor{
		repository: jobRepository,
		scheduler:  scheduler,
		client: &http.Client{
			Timeout: 30 * time.Second,
		},
		consumer: kafkaConsumer,
	}
}

func (e *JobExecutor) Start(ctx context.Context) error {
	return e.consumer.Consume(ctx, "job-events", e.handleJobEvent)
}

func (e *JobExecutor) handleJobEvent(event *events.JobEvent) error {
	switch event.Type {
	case events.JobCreated:
		return e.scheduleJob(event)
	case events.JobDeleted:
		return e.cancelJob(event.JobID)
	default:
		log.Printf("[WARN] Unknown event type: %s", event.Type)
		return nil
	}
}

func (e *JobExecutor) scheduleJob(event *events.JobEvent) error {
	jobData, err := e.repository.FindByID(context.Background(), event.JobID)
	if err != nil {
		return fmt.Errorf("failed to get job from repository: %w", err)
	}

	cronJob, err := e.scheduler.Cron(jobData.CronExpression).Tag(event.JobID).Do(e.executeJob, event.JobID)
	if err != nil {
		return fmt.Errorf("failed to schedule job: %w", err)
	}

	jobData.NextRun = cronJob.NextRun()
	if err := e.repository.Update(context.Background(), jobData); err != nil {
		return fmt.Errorf("failed to update job in repository: %w", err)
	}

	log.Printf("[INFO] Job %s scheduled for execution", event.JobID)
	return nil
}

func (e *JobExecutor) cancelJob(jobID string) error {
	e.scheduler.RemoveByTag(jobID)
	log.Printf("[INFO] Job %s cancelled", jobID)
	return nil
}


func (e *JobExecutor) executeJob(jobID string) {
	log.Printf("[INFO] Executing job %s", jobID)

	job, err := e.repository.FindByID(context.Background(), jobID)
	if err != nil {
		log.Printf("[ERROR] Failed to get job %s from repository: %v", jobID, err)
		return
	}

	startTime := time.Now()
	resp, err := e.client.Get(job.URL)

	job.LastRun = startTime
	if err != nil {
		job.LastStatus = 0
		log.Printf("[ERROR] Job %s failed: %v", jobID, err)
	} else {
		job.LastStatus = int32(resp.StatusCode)
		resp.Body.Close()
		log.Printf("[INFO] Job %s completed with status %d", jobID, resp.StatusCode)
	}

	if updateErr := e.repository.Update(context.Background(), job); updateErr != nil {
		log.Printf("[ERROR] Failed to update job %s in repository: %v", jobID, updateErr)
	}
}

func (e *JobExecutor) Stop() {
	e.scheduler.Stop()
	e.consumer.Close()
}
