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
	Conclusion   string     `gorm:"column:conclusion" json:"conclusion"`
	StartedAt    *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt   *time.Time `gorm:"column:finished_at" json:"finished_at"`
}

const (
	STATUS_QUEUED      = "QUEUED"
	STATUS_IN_PROGRESS = "IN_PROGRESS"
	STATUS_COMPLETED   = "COMPLETED"
)

// 新たに開始されたジョブか判定する
// QUEUED→IN_PROGRESS
func (after *Job) IsNewStarted(before *Job) bool {
	if before.JobId != after.JobId {
		return false
	}
	if before.Status != STATUS_QUEUED {
		return false
	}
	if after.Status != STATUS_IN_PROGRESS {
		return false
	}
	return true
}

// 新たに終了したジョブか判定する
// QUEUED→COMPLETED, IN_PROGRESS→COMPLETED
func (after *Job) IsNewFinished(before *Job) bool {
	if before.JobId != after.JobId {
		return false
	}
	if before.Status == STATUS_COMPLETED {
		return false
	}
	if after.Status != STATUS_COMPLETED {
		return false
	}
	return true
}
