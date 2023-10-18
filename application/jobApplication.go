package application

import (
	"log"
	"time"

	"github.com/kazuo278/dashboard/domain/model/github"
	"github.com/kazuo278/dashboard/domain/repository"
	"github.com/kazuo278/dashboard/domain/service"
)

type JobApplication interface {
	SetUpRunner(repositoryId string, repositoryName string, runId string, workflowRef string, jobName string, runAttempt string) error
	CompletedRunner(repositoryId string, runId string, jobName string, runAttempt string) error
}

type jobApplicationImpl struct {
	analyzer      analyzer
	jobRepository repository.JobRepository
	jobApi        service.JobApi
}

func NewJobAppication(jobRepositroy repository.JobRepository, jobApi service.JobApi) JobApplication {
	jobApplication := new(jobApplicationImpl)
	jobApplication.jobRepository = jobRepositroy
	jobApplication.jobApi = jobApi
	jobApplication.analyzer = newAnalyzer(jobRepositroy, jobApi)
	return jobApplication
}

// Set Up Runnerフェーズで実行するロジック
func (app *jobApplicationImpl) SetUpRunner(repositoryId string, repositoryName string, runId string, workflowRef string, jobName string, runAttempt string) error {
	log.Printf("INFO: [開始登録][Starting][RunID=%s][RunAttempt=%s][JobName=%s]ジョブ開始情報登録が開始されました", runId, runAttempt, jobName)
	// リポジトリを取得
	repository := app.jobRepository.GetRepositoryById(repositoryId)
	// リポジトリが存在しない(新規Actions実行したリポジトリである)場合
	if (*repository == github.Repository{}) {
		// リポジトリテーブルに登録
		log.Printf("INFO: [開始登録][InsertRepository][RunID=%s][RunAttempt=%s][JobName=%s]未登録リポジトリ%sをデータベースに登録します", runId, runAttempt, jobName, repositoryName)
		repository := new(github.Repository)
		repository.RepositoryId = repositoryId
		repository.RepositoryName = repositoryName
		app.jobRepository.CreateRepository(repository)
	}
	// ジョブリストの取得
	jobs, err := app.jobApi.GetJobList(runId, repositoryName, runAttempt)
	if err != nil {
		return err
	}
	for _, newJob := range *jobs {
		oldJob := app.jobRepository.GetJobById(newJob.JobId)
		// 新規ジョブの場合
		if (*oldJob == github.Job{}) {
			log.Printf("INFO: [開始登録][InsertJob][RunID=%s][RunAttempt=%s][JobName=%s][JobID=%s]新規ジョブをDBに登録します", runId, runAttempt, jobName, newJob.JobId)
			// 新規登録
			// APIでは連携されない情報をセット
			newJob.RunAttempt = runAttempt
			newJob.RepositoryId = repositoryId
			newJob.WorkflowRef = workflowRef
			app.jobRepository.CreateJob(&newJob)
		} else if newJob.IsNewStarted(oldJob) {
			// 新たに開始されたジョブの場合
			log.Printf("INFO: [開始登録][UpdateJob][RunID=%s][RunAttempt=%s][JobName=%s][JobID=%s]新たに開始状態となったジョブを更新します", runId, runAttempt, jobName, newJob.JobId)
			// APIでは連携されない情報をセット
			newJob.RunAttempt = runAttempt
			newJob.RepositoryId = repositoryId
			newJob.WorkflowRef = workflowRef
			app.jobRepository.UpdateJob(&newJob)
		}
	}
	log.Printf("INFO: [開始登録][Finished][RunID=%s][RunAttempt=%s][JobName=%s]ジョブ開始情報登録を終了します", runId, runAttempt, jobName)
	return nil
}

// Completed Runner フェーズで実行する関数
func (app *jobApplicationImpl) CompletedRunner(repositoryId string, runId string, jobName string, runAttempt string) error {
	log.Printf("INFO: [終了登録][Starting][RunID=%s][RunAttempt=%s][JobName=%s]ジョブ終了情報登録が開始しました", runId, runAttempt, jobName)
	go app.async(repositoryId, runId, jobName, runAttempt)
	log.Printf("INFO: [終了登録][Finishied][RunID=%s][RunAttempt=%s][JobName=%s]ジョブ終了情報登録を終了します", runId, runAttempt, jobName)
	return nil
}

func (app *jobApplicationImpl) async(repositoryId string, runId string, jobName string, runAttempt string) error {
	// リポジトリを取得
	repository := app.jobRepository.GetRepositoryById(repositoryId)

	// 対象ジョブステータスが終了するか、３回実行するまでジョブ情報を取得
	// ジョブ終了待ち基本時間(s)
	waitTimeSec := 5
	// 終了ジョブリスト
	finishedJobList := []github.Job{}
	for tryNum := 1; tryNum <= 3; tryNum++ {
		// ジョブ終了前はログ取得不可であるため、ジョブ終了まで待つ。
		time.Sleep(time.Duration(waitTimeSec*tryNum) * time.Second)
		// APIからジョブを取得
		newJobs, err := app.jobApi.GetJobList(runId, repository.RepositoryName, runAttempt)
		if err != nil {
			return err
		}
		for _, newJob := range *newJobs {
			oldJob := app.jobRepository.GetJobById(newJob.JobId)
			// 新たに終了したjobの場合
			if newJob.IsNewFinished(oldJob) {
				// APIでは連携されない情報をセット
				newJob.RunAttempt = runAttempt
				newJob.RepositoryId = repositoryId
				newJob.WorkflowRef = oldJob.WorkflowRef
				// 終了状態として更新する対象リストに追加
				finishedJobList = append(finishedJobList, newJob)
			}
		}
		// 新たな終了ジョブが存在する場合、ループを抜ける
		if len(finishedJobList) != 0 {
			break
		}
	}
	if len(finishedJobList) == 0 {
		log.Printf("INFO: [終了登録][RunID=%s][RunAttempt=%s][JobName=%s]新たに終了したジョブはありません", runId, runAttempt, jobName)
	}
	for _, finishedJob := range finishedJobList {
		// 終了ジョブを更新
		log.Printf("INFO: [終了登録][UpdateJob][RunID=%s][RunAttempt=%s][JobName=%s][JobID=%s]新たに終了状態となったジョブを更新します", runId, runAttempt, jobName, finishedJob.JobId)
		app.jobRepository.UpdateJob(&finishedJob)
		// ジョブ解析を実施
		app.analyzer.analyze(finishedJob.JobId, repository.RepositoryName)
	}
	return nil
}
