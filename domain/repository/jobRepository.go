package repository

import (
	"time"

	"github.com/kazuo278/dashboard/domain/model/github"
)

type JobRepository interface {
	// RepositoryのID検索
	GetRepositoryById(repositoryId string) *github.Repository
	// リポジトリの登録
	CreateRepository(repository *github.Repository) *github.Repository
	// ID検索
	GetJobById(jobId string) *github.Job
	// ID検索
	GetUnfinishedJob(from time.Time, to time.Time) *github.Job
	// ジョブの登録
	CreateJob(job *github.Job) *github.Job
	// ジョブの更新
	UpdateJob(job *github.Job) *github.Job
	// ジョブ詳細の登録
	CreateJobDetail(jobDetail *github.JobDetail) *github.JobDetail
}
