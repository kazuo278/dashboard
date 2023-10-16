package application

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/kazuo278/dashboard/domain/model/github"
	"github.com/kazuo278/dashboard/domain/repository"
	"github.com/kazuo278/dashboard/domain/service"
)

type analyzer interface {
	analyze(jobId string, repositoryName string)
}

type analyzerImpl struct {
	jobRepository repository.JobRepository
	jobApi        service.JobApi
}

func NewAnalyzer(jobRepositroy repository.JobRepository, jobApi service.JobApi) analyzer {
	analyzer := new(analyzerImpl)
	analyzer.jobRepository = jobRepositroy
	analyzer.jobApi = jobApi
	return analyzer
}

// ジョブ解析
func (analyzer analyzerImpl) analyze(jobId string, repositoryName string) {
	log.Printf("INFO: [ジョブ解析][Starting][JobID=%s]ログ解析を開始します", jobId)
	// ジョブ詳細をAPIから取得
	completedJob, err := analyzer.getJobFromApi(jobId, repositoryName)
	if err != nil || completedJob == nil {
		log.Printf("ERROR: [ジョブ解析][JobID=%s]GitHub APIからジョブの取得に失敗しました", jobId)
		return
	}
	// ログ取得
	jobLog, err := analyzer.jobApi.GetJobLog(jobId, repositoryName)
	if err != nil || jobLog == nil {
		log.Printf("ERROR: [ジョブ解析][JobID=%s]GitHub APIからジョブログの取得に失敗しました", jobId)
		return
	}
	// ログからActionとReusableワークフローの呼び出し情報を抽出
	jobDetailList := analyzer.extractJobDetails(*jobLog, jobId)
	log.Printf("INFO: [ジョブ解析][JobID=%s]合計%d件のActionまたはReusableWorkflowの呼び出しを見つけました", jobId, len(*jobDetailList))
	for _, detail := range *jobDetailList {
		// ジョブリストに登録
		log.Printf("INFO: [ジョブ解析][InsertJobDetail][JobID=%s][Type=%s]新規ジョブ詳細をDBに登録します", jobId, detail.Type)
		analyzer.jobRepository.CreateJobDetail(&detail)
	}
	log.Printf("INFO: [ジョブ解析][Finished][JobID=%s]ログ解析が完了しました", jobId)
}

// ジョブ取得関数
func (analyzer analyzerImpl) getJobFromApi(jobId string, repositoryName string) (*github.Job, error) {
	// APIからジョブを取得
	job, err := analyzer.jobApi.GetJob(jobId, repositoryName)
	if err != nil {
		return nil, err
	}
	if strings.ToUpper(job.Status) == github.STATUS_COMPLETED {
		// 取得完了したら、リトライを待たずループを抜ける
		log.Printf("INFO: [ジョブ解析][JobID=%s]GitHub APIからジョブを取得しました", jobId)
		return job, nil
	}
	return nil, fmt.Errorf("GitHub APIから終了したJobの取得に失敗しました")
}

// ログからジョブ詳細を取得
func (analyzer analyzerImpl) extractJobDetails(logByte []byte, jobId string) *[]github.JobDetail {
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
	jobDetailList = append(jobDetailList, *analyzer.generalizedExtractJobDetails(logByte, regexCallReusableWorkflowLog, jobId, github.TYPE_REUSABLE_WORKFLOW)...)
	jobDetailList = append(jobDetailList, *analyzer.generalizedExtractJobDetails(logByte, regexCallActionLogStr, jobId, github.TYPE_ACTION)...)
	return &jobDetailList
}

// ログから指定した正規表現のパターンマッチを行い、ジョブ取得を返却する
func (analyzer analyzerImpl) generalizedExtractJobDetails(logByte []byte, regex *regexp.Regexp, jobId string, jobDetailType string) *[]github.JobDetail {
	matches := regex.FindAllSubmatch(logByte, -1)
	jobDetailList := []github.JobDetail{}
	log.Printf("INFO: [ジョブ解析][JobID=%s]type=%sの呼び出しを%d件見つけました", jobId, jobDetailType, len(matches))
	for _, v := range matches {
		jobDetail := new(github.JobDetail)
		jobDetail.UsingPath = string(v[1]) + "/" + string(v[2])
		jobDetail.UsingRef = string(v[3])
		jobDetail.JobId = jobId
		jobDetail.Type = jobDetailType
		jobDetailList = append(jobDetailList, *jobDetail)
	}
	return &jobDetailList
}
