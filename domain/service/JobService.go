package service

import "github.com/kazuo278/action-dashboard/domain/model/github"

type JobService interface {
	// JobIdリストを取得する
	GetJobIdList(runId string, repositoryName string, runAttempt string) (*[]github.Job, error)
	// ジョブリストの取得
	GetJob(jobId string, repositoryName string) (*github.Job, error)
	// JobIdに紐づくログを取得する
	GetJobLog(jobId string, repositoryName string) (*[]byte, error)
	// ジョブを更新する
	UpdateJob(job *github.Job)
}
