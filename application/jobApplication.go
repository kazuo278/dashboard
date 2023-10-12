package application

import (
	"fmt"
	"log"
	"regexp"
	"strings"
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
func (app *jobApplicationImpl) SetUpRunner(repositoryId string, repositoryName string, runId string, workflowRef string, jobName string, runAttempt string) error {
	log.Printf("INFO: [RunID=%s][RunAttempt=%s][JobName=%s]ジョブが実行されました", runId, runAttempt, jobName)
	// リポジトリを取得
	repository := app.jobRepository.GetRepositoryById(repositoryId)
	// リポジトリが存在しない(新規Actions実行したリポジトリである)場合
	if (*repository == github.Repository{}) {
		// リポジトリテーブルに登録
		log.Printf("INFO: [RunID=%s][RunAttempt=%s][JobName=%s]リポジトリ名%sはデータベースに未登録のため登録します", runId, runAttempt, jobName, repositoryName)
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
			log.Printf("INFO: [JobID=%s]新規ジョブを登録します", newJob.JobId)
			// 新規登録
			// APIでは連携されない情報をセット
			newJob.RunAttempt = runAttempt
			newJob.RepositoryId = repositoryId
			newJob.WorkflowRef = workflowRef
			app.jobRepository.CreateJob(&newJob)
		} else if newJob.IsNewStarted(oldJob) {
			// 新たに開始されたジョブの場合
			log.Printf("INFO: [JobID=%s]新たに開始状態となったジョブを更新します", newJob.JobId)
			// APIでは連携されない情報をセット
			newJob.RunAttempt = runAttempt
			newJob.RepositoryId = repositoryId
			newJob.WorkflowRef = workflowRef
			app.jobRepository.UpdateJob(&newJob)
		}
	}
	log.Printf("INFO: [RunID=%s][RunAttempt=%s][JobName=%s]ジョブ開始情報登録を終了します", runId, runAttempt, jobName)
	return nil
}

// Completed Runner フェーズで実行する関数
func (app *jobApplicationImpl) CompletedRunner(repositoryId string, runId string, jobName string, runAttempt string) error {
	log.Printf("INFO: [RunID=%s][RunAttempt=%s][JobName=%s]ジョブが終了しました", runId, runAttempt, jobName)
	go app.async(repositoryId, runId, jobName, runAttempt)
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
			// TODO: ここで判定がおかしい
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
		log.Printf("INFO: [RunID=%s][RunAttempt=%s][JobName=%s]新たに終了したジョブはありません", runId, runAttempt, jobName)
	}
	for _, finishedJob := range finishedJobList {
		// 終了ジョブを更新
		app.jobRepository.UpdateJob(&finishedJob)
		// ジョブ解析を実施
		app.analyze(finishedJob.JobId, repository.RepositoryName)
	}
	return nil
}

// ジョブ解析
func (app jobApplicationImpl) analyze(jobId string, repositoryName string) {

	// ジョブ詳細をAPIから取得
	completedJob, err := app.getJobFromApi(jobId, repositoryName)
	if err != nil || completedJob == nil {
		log.Printf("ERROR: [JobID=%s] GitHub APIからジョブの取得に失敗しました", jobId)
		return
	}
	// ログ取得
	jobLog, err := app.jobApi.GetJobLog(jobId, repositoryName)
	if err != nil || jobLog == nil {
		log.Printf("ERROR: [JobID=%s] GitHub APIからジョブログの取得に失敗しました", jobId)
		return
	}
	log.Printf("INFO: [JobID=%s] ログ解析を実施します", jobId)
	// ログからActionとReusableワークフローの呼び出し情報を抽出
	jobDetailList := app.extractJobDetails(*jobLog, jobId)
	log.Printf("INFO:[JobID=%s] 合計%d件のActionまたはReusableWorkflowの呼び出しを見つけました", jobId, len(*jobDetailList))
	for _, detail := range *jobDetailList {
		// ジョブリストに登録
		// BUG 余計な空のデータを登録してしまう。
		app.jobRepository.CreateJobDetail(&detail)
	}
	log.Printf("INFO: [JobID=%s] ログ解析が完了しました", jobId)
}

// ジョブ取得関数
func (app jobApplicationImpl) getJobFromApi(jobId string, repositoryName string) (*github.Job, error) {
	// APIからジョブを取得
	job, err := app.jobApi.GetJob(jobId, repositoryName)
	if err != nil {
		return nil, err
	}
	if strings.ToUpper(job.Status) == github.STATUS_COMPLETED {
		// 取得完了したら、リトライを待たずループを抜ける
		log.Printf("INFO: [JobID=%s] GitHub APIからジョブを取得しました", jobId)
		return job, nil
	}
	return nil, fmt.Errorf("GitHub APIから終了したJobの取得に失敗しました")
}

// ログからジョブ詳細を取得
func (app jobApplicationImpl) extractJobDetails(logByte []byte, jobId string) *[]github.JobDetail {
	// Reusableワークフロー呼び出し判定用正規表現 $1=Organization, $2=Repository, $3=ref
	const REGEX_CALL_REUSABLE_WORKFLOW_LOG_STR = `(?m)^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{7}Z Uses: ([^/]+)/([^/]+)/[^@]+@([^ ]+) \([a-z0-9]+\)$`
	// Action呼び出し判定用正規表現 $1=Organization, $2=Repository, $3=ref
	const REGEX_CALL_ACTION_LOG_STR = `(?m)^\d{4}-\d{2}-\d{2}T\d{2}:\d{2}:\d{2}\.\d{7}Z Download action repository '([^/]+)/([^@]+)@(.+)' \(SHA:[a-z0-9]+\)$`
	// 戻り値の初期値定義
	jobDetailList := []github.JobDetail{}
	// Reusableワークフロー呼び出し判定用正規表現のコンパイル
	regexCallReusableWorkflowLog := regexp.MustCompile(REGEX_CALL_REUSABLE_WORKFLOW_LOG_STR)
	// Action呼び出し判定用正規表現のコンパイル
	regexCallActionLogStr := regexp.MustCompile(REGEX_CALL_ACTION_LOG_STR)
	// 抽出
	jobDetailList = append(jobDetailList, *app.generalizedExtractJobDetails(logByte, regexCallReusableWorkflowLog, jobId, github.TYPE_REUSABLE_WORKFLOW)...)
	jobDetailList = append(jobDetailList, *app.generalizedExtractJobDetails(logByte, regexCallActionLogStr, jobId, github.TYPE_ACTION)...)
	return &jobDetailList
}

// ログから指定した正規表現のパターンマッチを行い、ジョブ取得を返却する
func (app jobApplicationImpl) generalizedExtractJobDetails(logByte []byte, regex *regexp.Regexp, jobId string, jobDetailType string) *[]github.JobDetail {
	matches := regex.FindAllSubmatch(logByte, -1)
	jobDetailList := []github.JobDetail{}
	log.Printf("INFO: [JobID=%s] type=%sの呼び出しを%d件見つけました", jobId, jobDetailType, len(matches))
	for _, v := range matches {
		jobDetail := new(github.JobDetail)
		//　repositoryName="organization/repository"
		jobDetail.UsingPath = string(v[1]) + "/" + string(v[2])
		jobDetail.UsingRef = string(v[3])
		jobDetail.JobId = jobId
		jobDetail.Type = jobDetailType
		jobDetailList = append(jobDetailList, *jobDetail)
	}
	return &jobDetailList
}
