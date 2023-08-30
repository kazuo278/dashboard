package repository

import "github.com/kazuo278/dashboard/domain/model/github"

type JobRepository interface {
	// RepositoryのID検索
	GetRepositoryById(repositoryId string) *github.Repository
	// リポジトリの登録
	CreateRepository(repository *github.Repository) *github.Repository
	// ID代替の複数条件によるジョブの１件取得
	GetJobByIds(repositoryId string, runId string, jobName string, runAttempt string) *github.Job
	// ジョブの登録
	CreateJob(job *github.Job) *github.Job
	// ジョブの更新
	UpdateJob(job *github.Job) *github.Job
	// ジョブ詳細の登録
	CreateJobDetail(jobDetail *github.JobDetail) *github.JobDetail
}
