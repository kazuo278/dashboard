package application

import (
	"github.com/kazuo278/dashboard/application/dto"
	"github.com/kazuo278/dashboard/application/repository"
)

type DashboardApplication interface {
	GetJobs(limit int, offset int, jobId string, repositoryId string, repositoryName string, workflowRef string, jobName string, runId string, runAttempt string, status string, startedAt string, finishedAt string) (*[]dto.RepoJob, int)
	GetJobCount(repositoryName string, startedAt string, finishedAt string) (*[]dto.RepoCount, int)
	GetJobTime(repositoryName string, startedAt string, finishedAt string) (*[]dto.RepoTime, int)
	GetJobDetails(limit int, offset int, jobId string, repositoryId string, repositoryName string, usingPath string, usingRef string, jobName string, runId string, runAttempt string, typeName string, startedAt string, finishedAt string) (*[]dto.RepoJobDetail, int)
}

type dashboardApplicationImpl struct {
	repository repository.DashboardRepository
}

func NewDashboardApplication(dashboardRepository repository.DashboardRepository) DashboardApplication {
	dashboardApplication := new(dashboardApplicationImpl)
	dashboardApplication.repository = dashboardRepository
	return dashboardApplication
}

// 実行履歴を取得する
func (app dashboardApplicationImpl) GetJobs(limit int, offset int, jobId string, repositoryId string, repositoryName string, workflowRef string, jobName string, runId string, runAttempt string, status string, startedAt string, finishedAt string) (*[]dto.RepoJob, int) {
	jobs, count := app.repository.GetJobs(limit, offset, jobId, repositoryId, repositoryName, workflowRef, jobName, runId, runAttempt, status, startedAt, finishedAt)
	return jobs, count
}

// リポジトリごとの実行回数を取得する
func (app dashboardApplicationImpl) GetJobCount(repositoryName string, startedAt string, finishedAt string) (*[]dto.RepoCount, int) {
	result, count := app.repository.GetJobCount(repositoryName, startedAt, finishedAt)
	return result, count
}

// リポジトリごとの実行時間を取得する
func (app dashboardApplicationImpl) GetJobTime(repositoryName string, startedAt string, finishedAt string) (*[]dto.RepoTime, int) {
	result, totalSeconds := app.repository.GetJobTime(repositoryName, startedAt, finishedAt)
	return result, totalSeconds
}

// ジョブ詳細を取得する
func (app dashboardApplicationImpl) GetJobDetails(limit int, offset int, jobId string, repositoryId string, repositoryName string, usingPath string, usingRef string, jobName string, runId string, runAttempt string, typeName string, startedAt string, finishedAt string) (*[]dto.RepoJobDetail, int) {
	result, count := app.repository.GetJobDetails(limit, offset, jobId, repositoryId, repositoryName, usingPath, usingRef, jobName, runId, runAttempt, typeName, startedAt, finishedAt)
	return result, count
}
