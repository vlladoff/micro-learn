package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-redis/redis"
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

type JobRepository interface {
	Save(ctx context.Context, job *Job) error
	FindByID(ctx context.Context, id string) (*Job, error)
	Delete(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	Update(ctx context.Context, job *Job) error
}

type RedisJobRepository struct {
	client *redis.Client
}

func NewRedisJobRepository(client *redis.Client) JobRepository {
	return &RedisJobRepository{
		client: client,
	}
}

func (r *RedisJobRepository) Save(ctx context.Context, job *Job) error {
	jobJSON, err := json.Marshal(job)
	if err != nil {
		return fmt.Errorf("failed to marshal job: %w", err)
	}

	key := fmt.Sprintf("job:%s", job.ID)
	if err := r.client.Set(key, jobJSON, 0).Err(); err != nil {
		return fmt.Errorf("failed to save job to redis: %w", err)
	}

	return nil
}

func (r *RedisJobRepository) FindByID(ctx context.Context, id string) (*Job, error) {
	key := fmt.Sprintf("job:%s", id)
	jobJSON, err := r.client.Get(key).Result()
	if err != nil {
		if err == redis.Nil {
			return nil, fmt.Errorf("job not found")
		}
		return nil, fmt.Errorf("failed to get job from redis: %w", err)
	}

	var job Job
	if err := json.Unmarshal([]byte(jobJSON), &job); err != nil {
		return nil, fmt.Errorf("failed to unmarshal job: %w", err)
	}

	return &job, nil
}

func (r *RedisJobRepository) Delete(ctx context.Context, id string) error {
	key := fmt.Sprintf("job:%s", id)
	result, err := r.client.Del(key).Result()
	if err != nil {
		return fmt.Errorf("failed to delete job from redis: %w", err)
	}
	if result == 0 {
		return fmt.Errorf("job not found")
	}
	return nil
}

func (r *RedisJobRepository) Exists(ctx context.Context, id string) (bool, error) {
	key := fmt.Sprintf("job:%s", id)
	result, err := r.client.Exists(key).Result()
	if err != nil {
		return false, fmt.Errorf("failed to check job existence: %w", err)
	}
	return result > 0, nil
}

func (r *RedisJobRepository) Update(ctx context.Context, job *Job) error {
	exists, err := r.Exists(ctx, job.ID)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("job not found")
	}

	return r.Save(ctx, job)
}
