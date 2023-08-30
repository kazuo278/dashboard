package dto

import (
	"time"
)

type RepoJob struct {
	JobId          string     `gorm:"primaryKey;column:job_id" json:"job_id"`
	RepositoryId   string     `gorm:"primaryKey;column:repository_id" json:"repository_id"`
	RepositoryName string     `gorm:"column:repository_name" json:"repository_name"`
	RunId          string     `gorm:"primaryKey;column:run_id" json:"run_id"`
	WorkflowRef    string     `gorm:"column:workflow_ref" json:"workflow_ref"`
	JobName        string     `gorm:"primaryKey;column:job_name" json:"job_name"`
	RunAttempt     string     `gorm:"primaryKey;column:run_attempt" json:"run_attempt"`
	Status         string     `gorm:"column:status" json:"status"`
	StartedAt      *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt     *time.Time `gorm:"column:finished_at" json:"finished_at"`
}
