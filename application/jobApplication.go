package application

import (
	"log"
	"time"

	"github.com/kazuo278/dashboard/domain/model/github"
	"github.com/kazuo278/dashboard/domain/repository"
	"github.com/kazuo278/dashboard/domain/service"
)

type JobApplication interface {
	SetUpRunner(repositoryId string, repositoryName string, runId string, workflowRef string, jobName string, runAttempt string) (*github.Job, error)
	CompletedRunner(repositoryId string, runId string, jobName string, runAttempt string) *github.Job
}

type jobApplicationImpl struct {
	jobRepository repository.JobRepository
	jobApi        service.JobApi
}

func NewJobAppication(jobRepositroy repository.JobRepository, jobApi service.JobApi) JobApplication {
	jobApplication := new(jobApplicationImpl)
	jobApplication.jobRepository = jobRepositroy
	jobApplication.jobApi = jobApi
	return jobApplication
}

// Set Up Runnerフェーズで実行するロジック
func (app *jobApplicationImpl) SetUpRunner(repositoryId string, repositoryName string, runId string, workflowRef string, jobName string, runAttempt string) (*github.Job, error) {
	// 現在実行中のジョブがdbに登録済みか調べる
	currentJob := app.jobRepository.GetJobByIds(repositoryId, runId, jobName, runAttempt)
	log.Printf("INFO: [JobId=%s]ジョブ開始情報登録を開始します", currentJob.JobId)
	// リポジトリを取得
	repository := app.jobRepository.GetRepositoryById(repositoryId)
	// リポジトリが存在しない(新規Actions実行したリポジトリである)場合
	if (*repository == github.Repository{}) {
		// リポジトリテーブルに登録
		log.Printf("INFO: [JobId=%s]リポジトリ名%sはデータベースに未登録のため登録します", currentJob.RunId, repositoryName)
		repository := new(github.Repository)
		repository.RepositoryId = repositoryId
		repository.RepositoryName = repositoryName
		app.jobRepository.CreateRepository(repository)
	}
	// 未登録の場合、RunId、RunAttemptに紐づく全てのジョブをdbに登録する
	if (*currentJob == github.Job{}) {
		log.Printf("INFO: [JobId=%s]本ジョブはデータベースに未登録のため登録します", currentJob.JobId)
		// ジョブリストの取得
		jobs, err := app.jobApi.GetJobList(runId, repositoryName, runAttempt)
		if err != nil {
			return nil, err
		}
		for _, job := range *jobs {
			// 未登録のジョブ情報を登録
			job.RepositoryId = repositoryId
			job.WorkflowRef = workflowRef
			job.RunAttempt = runAttempt
			if job.JobName == jobName {
				// 現在実行中のジョブは、Startedステータスとしてdbに登録
				job.Status = "STARTED"
				job.StartedAt = nowJST()
				job.FinishedAt = nil
				// 戻り値用変数に詰め替え
				currentJob = &job
			} else {
				// 現在実行中のジョブ以外は、Queuedステータスとしてdbに登録
				job.Status = "QUEUED"
				job.StartedAt = nil
				job.FinishedAt = nil
			}
			app.jobRepository.CreateJob(&job)
		}
	} else {
		// 登録済みの場合
		// 実行中のジョブを実行中ステータスに変更
		currentJob.StartedAt = nowJST()
		currentJob.Status = "STARTED"
		app.jobRepository.UpdateJob(currentJob)
	}
	log.Printf("INFO: [JobId=%s]ジョブ開始情報登録を終了します", currentJob.JobId)
	return currentJob, nil
}

// Completed Runner フェーズで実行する関数
func (app *jobApplicationImpl) CompletedRunner(repositoryId string, runId string, jobName string, runAttempt string) *github.Job {
	// 実行中のジョブをdbから取得
	job := app.jobRepository.GetJobByIds(repositoryId, runId, jobName, runAttempt)
	if job.FinishedAt == nil && job.Status != "COMPLETED" {
		log.Printf("INFO: [JobId=%s]ジョブ終了情報登録を開始します", job.JobId)
		// 終了ステータスに変更
		job.Status = "COMPLETED"
		// 終了時刻をセット
		job.FinishedAt = nowJST()
		// DBを更新
		app.jobRepository.UpdateJob(job)
		log.Printf("INFO: [JobId=%s]ジョブ終了情報登録を終了します", job.JobId)
	}
	return job
}

func nowJST() *time.Time {
	jst, err := time.LoadLocation("Asia/Tokyo")
	if err != nil {
		panic(err)
	}
	nowJST := time.Now().In(jst)
	return &nowJST
}
