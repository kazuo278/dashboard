package dto

import (
	"time"
)

type RepoJobDetail struct {
	JobDetailId    uint       `gorm:"primaryKey;column:job_detail_id" json:"job_detail_id"`
	JobId          string     `gorm:"primaryKey;column:job_id" json:"job_id"`
	RepositoryId   string     `gorm:"primaryKey;column:repository_id" json:"repository_id"`
	RepositoryName string     `gorm:"column:repository_name" json:"repository_name"`
	RunId          string     `gorm:"primaryKey;column:run_id" json:"run_id"`
	JobName        string     `gorm:"primaryKey;column:job_name" json:"job_name"`
	RunAttempt     string     `gorm:"primaryKey;column:run_attempt" json:"run_attempt"`
	Type           string     `gorm:"column:type" json:"type"`
	UsingPath      string     `gorm:"column:using_path" json:"using_path"`
	UsingRef       string     `gorm:"column:using_ref" json:"using_ref"`
	StartedAt      *time.Time `gorm:"column:started_at" json:"started_at"`
	FinishedAt     *time.Time `gorm:"column:finished_at" json:"finished_at"`
}
