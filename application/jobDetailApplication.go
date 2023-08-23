package application

import (
	"fmt"
	"regexp"
	"time"

	"github.com/kazuo278/action-dashboard/infrastruncture/database"
	"github.com/kazuo278/action-dashboard/infrastruncture/restapi"
	"github.com/kazuo278/action-dashboard/model/github"
)

func Analysis(jobId string, repositoryId string) (*github.Job, *[]github.JobDetail, error) {
	// リポジトリ取得
	repository := database.GetRepositoryById(repositoryId)
	// 終了ジョブの取得
	completedJob, err := getJob(jobId, repository.RepositoryName)
	if err != nil {
		return nil, nil, err
	}
	// ログの取得
	log, err := restapi.GetJobLog(jobId, repository.RepositoryName)
	// ジョブ詳細の取得
	jobDetailList := analysisLog(*log, jobId)
	// TODO：ジョブリストDB登録
	// 返却
	return completedJob, jobDetailList, nil
}

func getJob(jobId string, repositoryName string) (*github.Job, error) {
	// 対象ジョブステータスが終了するか、３回実行するまでジョブ情報を取得
	// ジョブ終了待ち基本時間(s)
	waitTimeSec := 5
	for tryNum := 1; tryNum <= 3; tryNum++ {
		// ジョブ終了前はログ取得不可であるため、ジョブ終了まで待つ。
		time.Sleep(time.Duration(waitTimeSec*tryNum) * time.Second)
		// ジョブを取得
		job, err := restapi.GetJob(jobId, repositoryName)
		if err != nil {
			return nil, err
		}
		if job.Status == "completed" {
			// 取得完了したら、リトライを待たずループを抜ける
			return job, nil
		}
	}
	return nil, fmt.Errorf("終了したJobの取得に失敗しました")
}

// ログからジョブ詳細を取得
func analysisLog(log []byte, jobId string) *[]github.JobDetail {
	// Reusableワークフロー呼び出し判定用正規表現 $1=Organization, $2=Repository, $3=ref
	const REGEX_CALL_REUSABLE_WORKFLOW_LOG_STR = `(?m)^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{7}Z Uses: ([^/]+)/([^/]+)/.+\\[^@]+@([^ ]+) \\([a-z0-9]+\\)$`
	// Action呼び出し判定用正規表現 $1=Organization, $2=Repository, $3=ref
	const REGEX_CALL_ACTION_LOG_STR = `(?m)^\\d{4}-\\d{2}-\\d{2}T\\d{2}:\\d{2}:\\d{2}\\.\\d{7}Z Download action repository '([^/]+)/([^@]+)@(.+)' \\(SHA:[a-z0-9]+\\)$`
	// 戻り値の初期値定義
	jobDetailList := []github.JobDetail{}
	// Reusableワークフロー呼び出し判定用正規表現のコンパイル
	regexCallReusableWorkflowLog := regexp.MustCompile(REGEX_CALL_REUSABLE_WORKFLOW_LOG_STR)
	// Action呼び出し判定用正規表現のコンパイル
	regexCallActionLogStr := regexp.MustCompile(REGEX_CALL_REUSABLE_WORKFLOW_LOG_STR)
	// 抽出
	jobDetailList = append(jobDetailList, *extractJobDetails(log, regexCallReusableWorkflowLog, jobId, "REUSABLE_WORKFLOW")...)
	jobDetailList = append(jobDetailList, *extractJobDetails(log, regexCallActionLogStr, jobId, "ACTION")...)
	return &jobDetailList
}

// ログから指定した正規表現のパターンマッチを行い、ジョブ取得を返却する
func extractJobDetails(log []byte, regex *regexp.Regexp, jobId string, jobDetailType string) *[]github.JobDetail {
	matches := regex.FindAllSubmatch(log, -1)
	jobDetailList := make([]github.JobDetail, len(matches))
	for _, v := range matches {
		jobDetail := new(github.JobDetail)
		//　repositoryName="organization/repository"
		jobDetail.UsingRepositoryName = string(v[1]) + "/" + string(v[2])
		jobDetail.UsingRepositoryRef = string(v[3])
		jobDetail.JobId = jobId
		jobDetail.Type = jobDetailType
		jobDetailList = append(jobDetailList, *jobDetail)
	}
	return &jobDetailList
}
