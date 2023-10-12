package repository

import "github.com/kazuo278/dashboard/application/dto"

type DashboardRepository interface {
	GetJobs(limit int, offset int, jobId string, repositoryId string, repositoryName string, workflowRef string, jobName string, runId string, runAttempt string, status string, conclusion string, startedAt string, finishedAt string) (*[]dto.RepoJob, int)
	GetJobCount(repositoryName string, startedAt string, finishedAt string) (*[]dto.RepoCount, int)
	GetJobTime(repositoryName string, startedAt string, finishedAt string) (*[]dto.RepoTime, int)
	GetJobDetails(limit int, offset int, jobId string, repositoryId string, repositoryName string, usingPath string, usingRef string, jobName string, runId string, runAttempt string, typeName string, startedAt string, finishedAt string) (*[]dto.RepoJobDetail, int)
}
