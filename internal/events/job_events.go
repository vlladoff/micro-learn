package events

import (
	"encoding/json"
	"time"
)

type EventType string

const (
	JobCreated  EventType = "job.created"
	JobUpdated  EventType = "job.updated"
	JobDeleted  EventType = "job.deleted"
	JobExecuted EventType = "job.executed"
)

type JobEvent struct {
	Type      EventType   `json:"type"`
	JobID     int64       `json:"job_id"`
	Timestamp time.Time   `json:"timestamp"`
	Data      interface{} `json:"data"`
}

type JobCreatedData struct {
	URL            string `json:"url"`
	CronExpression string `json:"cron_expression"`
	Comment        string `json:"comment"`
}

type JobExecutedData struct {
	Status   int32  `json:"status"`
	Response string `json:"response,omitempty"`
	Error    string `json:"error,omitempty"`
}

func (e *JobEvent) ToJSON() ([]byte, error) {
	return json.Marshal(e)
}

func (e *JobEvent) FromJSON(data []byte) error {
	return json.Unmarshal(data, e)
}
