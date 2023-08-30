package application

import (
	"fmt"
	"log"
	"regexp"
	"time"

	"github.com/kazuo278/dashboard/domain/model/github"
	"github.com/kazuo278/dashboard/domain/repository"
	"github.com/kazuo278/dashboard/domain/service"
)

type JobAnalyzerApplication interface {
	Analyze(jobId string, repositoryId string)
}

type jobAnalyzerApplicationImpl struct {
	jobRepository repository.JobRepository
	jobApi        service.JobApi
}

func NewJobAnalyzerAppication(jobRepositroy repository.JobRepository, jobApi service.JobApi) JobAnalyzerApplication {
	jobAnalyzerApplication := new(jobAnalyzerApplicationImpl)
	jobAnalyzerApplication.jobRepository = jobRepositroy
	jobAnalyzerApplication.jobApi = jobApi
	return jobAnalyzerApplication
}

func (app jobAnalyzerApplicationImpl) Analyze(jobId string, repositoryId string) {
	// リポジトリをdbから取得
	repository := app.jobRepository.GetRepositoryById(repositoryId)

	// 実行中のジョブの最新情報をAPIから取得
	completedJob, err := app.getJobFromApi(jobId, repository.RepositoryName)
	if err != nil || completedJob == nil {
		log.Printf("ERROR: [JobID=%s] GitHub APIからジョブの取得に失敗しました", jobId)
		return
	}
	// ログ取得
	jobLog, err := app.jobApi.GetJobLog(jobId, repository.RepositoryName)
	if err != nil || jobLog == nil {
		log.Printf("ERROR: [JobID=%s] GitHub APIからジョブログの取得に失敗しました", jobId)
		return
	}
	log.Printf("INFO: [JobID: %s] ログ解析を実施します", jobId)
	// ログからActionとReusableワークフローの呼び出し情報を抽出
	jobDetailList := app.extractJobDetails(*jobLog, jobId)
	log.Printf("INFO:[JobId: %s] 合計%d件のActionまたはReusableWorkflowの呼び出しを見つけました", jobId, len(*jobDetailList))
	for _, detail := range *jobDetailList {
		// ジョブリストに登録
		// BUG 余計な空のデータを登録してしまう。
		app.jobRepository.CreateJobDetail(&detail)
	}
	log.Printf("INFO: [JobID: %s] ログ解析が完了しました", jobId)
}

// ジョブ取得関数
func (app jobAnalyzerApplicationImpl) getJobFromApi(jobId string, repositoryName string) (*github.Job, error) {
	// 対象ジョブステータスが終了するか、３回実行するまでジョブ情報を取得
	// ジョブ終了待ち基本時間(s)
	waitTimeSec := 5
	for tryNum := 1; tryNum <= 3; tryNum++ {
		// ジョブ終了前はログ取得不可であるため、ジョブ終了まで待つ。
		time.Sleep(time.Duration(waitTimeSec*tryNum) * time.Second)
		// APIからジョブを取得
		job, err := app.jobApi.GetJob(jobId, repositoryName)
		if err != nil {
			return nil, err
		}
		if job.Status == "completed" {
			// 取得完了したら、リトライを待たずループを抜ける
			log.Printf("INFO: [JobID=%s] GitHub APIからジョブを取得しました", jobId)
			return job, nil
		}
	}
	// 3回リトライまでに終了を確認できない場合
	log.Printf("ERROR: [JobID=%s] GitHub APIからジョブログの取得に失敗しました", jobId)
	return nil, fmt.Errorf("GitHub APIから終了したJobの取得に失敗しました")
}

// ログからジョブ詳細を取得
func (app jobAnalyzerApplicationImpl) extractJobDetails(logByte []byte, jobId string) *[]github.JobDetail {
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
	jobDetailList = append(jobDetailList, *app.generalizedExtractJobDetails(logByte, regexCallReusableWorkflowLog, jobId, "REUSABLE_WORKFLOW")...)
	jobDetailList = append(jobDetailList, *app.generalizedExtractJobDetails(logByte, regexCallActionLogStr, jobId, "ACTION")...)
	return &jobDetailList
}

// ログから指定した正規表現のパターンマッチを行い、ジョブ取得を返却する
func (app jobAnalyzerApplicationImpl) generalizedExtractJobDetails(logByte []byte, regex *regexp.Regexp, jobId string, jobDetailType string) *[]github.JobDetail {
	matches := regex.FindAllSubmatch(logByte, -1)
	jobDetailList := []github.JobDetail{}
	log.Printf("INFO: [JobId=%s] type=%sの呼び出しを%d件見つけました", jobId, jobDetailType, len(matches))
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
