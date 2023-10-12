package dto

import "time"

type GitHubJob struct {
	JobId       int        `json:"id"`
	RunId       int        `json:"run_id"`
	JobName     string     `json:"name"`
	Status      string     `json:"status"`
	Conclusion  string     `json:"conclusion"`
	StartedAt   *time.Time `gorm:"column:started_at" json:"started_at"`
	CompletedAt *time.Time `gorm:"column:completed_at" json:"completed_at"`
}
