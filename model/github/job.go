package github

import (
	"time"
)

type Job struct {
	JobId        string     `gorm:"primaryKey;column:job_id" json:"job_id"`
	RunId        string     `gorm:"column:run_id" json:"run_id"`
	RunAttempt   string     `gorm:"column:run_attempt" json:"run_attempt"`
	RepositoryId string     `gorm:"column:repository_id" json:"repository_id"`
	WorkflowRef  string     `gorm:"column:workflow_ref" json:"workflow_ref"`
	JobName      string     `gorm:"column:job_name" json:"job_name"`
	Status       string     `gorm:"column:status" json:"status"`
	StartedAt    *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt   *time.Time `gorm:"column:finished_at" json:"finished_at"`
}
