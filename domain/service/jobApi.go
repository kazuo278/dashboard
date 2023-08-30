package service

import "github.com/kazuo278/dashboard/domain/model/github"

type JobApi interface {
	// JobIdリストを取得する
	GetJobList(runId string, repositoryName string, runAttempt string) (*[]github.Job, error)
	// ジョブリストの取得
	GetJob(jobId string, repositoryName string) (*github.Job, error)
	// JobIdに紐づくログを取得する
	GetJobLog(jobId string, repositoryName string) (*[]byte, error)
}
