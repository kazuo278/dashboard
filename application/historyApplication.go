package application

import (
	"log"
	"time"

	"github.com/kazuo278/action-dashboard/infrastruncture/database"
	"github.com/kazuo278/action-dashboard/infrastruncture/restapi"
	"github.com/kazuo278/action-dashboard/model/github"
)

// Set Up Runnerフェーズで実行するロジック
func SetUpRunner(repositoryId string, repositoryName string, runId string, workflowRef string, jobName string, runAttempt string) (*github.Job, error) {
	// リポジトリを取得
	repository := database.GetRepositoryById(repositoryId)
	// リポジトリが存在しない(新規Actions実行したリポジトリである)場合
	if (*repository == github.Repository{}) {
		// リポジトリテーブルに登録
		log.Print("リポジトリ名" + repositoryName + "は一覧に存在しないため登録します")
		repository := new(github.Repository)
		repository.RepositoryId = repositoryId
		repository.RepositoryName = repositoryName
		database.CreateRepository(repository)
	}
	// 現在実行中のジョブがdbに登録済みか調べる
	currentJob := database.GetJobByIds(repositoryId, runId, jobName, runAttempt)
	// 未登録の場合、RunId、RunAttemptに紐づく全てのジョブをdbに登録する
	if (*currentJob == github.Job{}) {
		// ジョブリストの取得
		jobs, err := restapi.GetJobList(runId, repositoryName, runAttempt)
		if err != nil {
			return nil, err
		}
		for _, job := range *jobs {
			// 未登録のジョブ情報を登録
			job.RepositoryId = repositoryId
			job.RunId = runId
			job.WorkflowRef = workflowRef
			job.JobName = jobName
			job.RunAttempt = runAttempt
			if job.JobName == jobName {
				// 現在実行中のジョブは、Startedステータスとしてdbに登録
				job.Status = "STARTED"
				job.StartedAt = nowJST()
				// 戻り値用変数に詰め替え
				currentJob = &job
			} else {
				// 現在実行中のジョブ以外は、Queuedステータスとしてdbに登録
				job.Status = "QUEUED"
			}
			database.CreateJob(&job)
		}
	} else {
		// 登録済みの場合
		// 実行中のジョブを実行中ステータスに変更
		currentJob.StartedAt = nowJST()
		currentJob.Status = "STARTED"
		database.UpdateJob(currentJob)
	}
	return currentJob, nil
}

// Completed Runner フェーズで実行する関数
func CompletedRunner(repositoryId string, runId string, jobName string, runAttempt string) *github.Job {
	// 更新対象を取得
	job := database.GetJobByIds(repositoryId, runId, jobName, runAttempt)
	// ステータスを変更
	job.FinishedAt = nowJST()
	job.Status = "FINISHED"
	// 更新
	database.UpdateJob(job)
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
