package restapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/kazuo278/dashboard/domain/model/github"
	"github.com/kazuo278/dashboard/domain/service"
	"github.com/kazuo278/dashboard/infrastruncture/restapi/dto"
)

const baseUrl string = "https://api.github.com/"
const versionHeader string = "X-GitHub-Api-Version"
const version string = "2022-11-28"
const AcceptHeader string = "Accept"
const accept string = "application/vnd.github+json"
const AuthorizationHeader string = "Authorization"

type gitHubJobApiImpl struct {
}

func NewGitHubJobApi() service.JobApi {
	api := new(gitHubJobApiImpl)
	return api
}

// ジョブリストの取得
func (api *gitHubJobApiImpl) GetJobList(runId string, repositoryName string, runAttempt string) (*[]github.Job, error) {
	// 参考：https://docs.github.com/ja/rest/actions/workflow-jobs?apiVersion=2022-11-28#list-jobs-for-a-workflow-run-attempt
	// URI：https://api.github.com/repos/OWNER/REPO/actions/runs/RUN_ID/attempts/ATTEMPT_NUMBER/jobs
	jobListEndpointUri := baseUrl + "repos/" + repositoryName + "/actions/runs/" + runId + "/attempts/" + runAttempt + "/jobs"
	log.Printf("INFO: [ジョブリスト取得API][RunID=%s][RunAttempt=%s]%sに接続します", runId, runAttempt, jobListEndpointUri)
	// トークン取得
	organizationToken, err := getOrganizationToken()
	if err != nil {
		return nil, err
	}

	// リクエスト送信
	request, _ := http.NewRequest(http.MethodGet, jobListEndpointUri, nil)
	request.Header.Set(versionHeader, version)
	request.Header.Set(AcceptHeader, accept)
	request.Header.Set(AuthorizationHeader, "Bearer "+organizationToken)
	client := new(http.Client)

	// レスポンス受信
	response, err := client.Do(request)
	if err != nil {
		log.Printf("ERROR: [ジョブリスト取得API][RunID=%s][RunAttempt=%s]ジョブリストの取得に失敗しました", runId, runAttempt)
		return nil, fmt.Errorf("error: エンドポイント接続時にエラーが発生しました: %w", err)
	}

	defer response.Body.Close()
	// 正常応答
	if response.StatusCode == 200 {
		// レスポンスのJSON文字列を構造体GitHubJobListにマッピング
		result := dto.GitHubJobList{}
		json.NewDecoder(response.Body).Decode(&result)
		log.Printf("INFO: [ジョブリスト取得API][RunID=%s][RunAttempt=%s]ジョブリストの取得に成功しました", runId, runAttempt)
		return convertGitHubJobListToJobList(result.Jobs), nil
		// エラー応答
	} else {
		result := dto.GitHubErrorResponse{}
		json.NewDecoder(response.Body).Decode(&result)
		log.Printf("ERROR: [ジョブリスト取得API][RunID=%s][RunAttempt=%s]ジョブリストの取得に失敗しました", runId, runAttempt)
		log.Printf("ERROR: [ジョブリスト取得API][RunID=%s][RunAttempt=%s]エラーコード: %s", runId, runAttempt, response.Status)
		log.Printf("ERROR: [ジョブリスト取得API][RunID=%s][RunAttempt=%s]レスポンスエラー: %s", runId, runAttempt, result.Message)
		log.Printf("ERROR: [ジョブリスト取得API][RunID=%s][RunAttempt=%s]ドキュメント: %s", runId, runAttempt, result.DocumentationUrl)
		return nil, fmt.Errorf("GitHub API Error : %s", result.Message)
	}
}

// ジョブの取得
func (api *gitHubJobApiImpl) GetJob(jobId string, repositoryName string) (*github.Job, error) {
	// 参考：https://docs.github.com/ja/rest/actions/workflow-jobs?apiVersion=2022-11-28#get-a-job-for-a-workflow-run
	// URI：https://api.github.com/repos/OWNER/REPO/actions/jobs/JOB_ID
	endpointUri := baseUrl + "repos/" + repositoryName + "/actions/jobs/" + jobId
	log.Printf("INFO: [ジョブ取得API][JobID=%s]%sに接続します", jobId, endpointUri)
	// トークン取得
	organizationToken, err := getOrganizationToken()
	if err != nil {
		return nil, err
	}

	// リクエスト送信
	request, _ := http.NewRequest(http.MethodGet, endpointUri, nil)
	request.Header.Set(versionHeader, version)
	request.Header.Set(AcceptHeader, accept)
	request.Header.Set(AuthorizationHeader, "Bearer "+organizationToken)
	client := new(http.Client)

	// レスポンス受信
	response, err := client.Do(request)
	if err != nil {
		log.Printf("ERROR: [ジョブ取得API][JobID=%s]ジョブの取得に失敗しました", jobId)
		return nil, fmt.Errorf("error: エンドポイント接続時にエラーが発生しました: %w", err)
	}

	defer response.Body.Close()
	// 正常応答
	if response.StatusCode == 200 {
		// レスポンスのJSON文字列を構造体GitHubResponseにマッピング
		result := github.Job{}
		json.NewDecoder(response.Body).Decode(&result)
		log.Printf("INFO: [ジョブ取得API][JobID=%s]ジョブ取得に成功しました", jobId)
		return &result, nil
		// エラー応答
	} else {
		result := dto.GitHubErrorResponse{}
		json.NewDecoder(response.Body).Decode(&result)
		log.Printf("ERROR: [ジョブ取得API][JobID=%s]ジョブの取得に失敗しました", jobId)
		log.Printf("ERROR: [ジョブ取得API][JobID=%s]エラーコード: %s", jobId, response.Status)
		log.Printf("ERROR: [ジョブ取得API][JobID=%s]レスポンスエラー: %s", jobId, result.Message)
		log.Printf("ERROR: [ジョブ取得API][JobID=%s]ドキュメント: %s", jobId, result.DocumentationUrl)
		return nil, fmt.Errorf("GitHub API Error : %s", result.Message)
	}
}

// JobIdに紐づくログを取得する
func (api *gitHubJobApiImpl) GetJobLog(jobId string, repositoryName string) (*[]byte, error) {
	// 参考：https://docs.github.com/ja/rest/actions/workflow-jobs?apiVersion=2022-11-28#download-job-logs-for-a-workflow-run
	// URI：https://api.github.com/repos/OWNER/REPO/actions/jobs/JOB_ID/logs
	endpointUri := baseUrl + "repos/" + repositoryName + "/actions/jobs/" + jobId + "/logs"
	log.Printf("INFO: [ログ取得API][JobID=%s]%sに接続します", jobId, endpointUri)
	// トークン取得
	organizationToken, err := getOrganizationToken()
	if err != nil {
		return nil, err
	}

	// リクエスト送信
	request, _ := http.NewRequest(http.MethodGet, endpointUri, nil)
	request.Header.Set(versionHeader, version)
	request.Header.Set(AcceptHeader, accept)
	request.Header.Set(AuthorizationHeader, "Bearer "+organizationToken)
	client := new(http.Client)

	// レスポンス受信
	// 本来302レスポンスを受けリダイレクトする必要があるが、net/httpの仕組みにより自動でリダイレクトする。
	response, err := client.Do(request)
	if err != nil {
		log.Printf("ERROR: [ログ取得API][JobID=%s]ログの取得に失敗しました", jobId)
		return nil, fmt.Errorf("error: エンドポイント接続時にエラーが発生しました: %w", err)
	}

	defer response.Body.Close()
	// 正常応答
	if response.StatusCode == 200 {
		// レスポンスのログ文字列を返却
		result, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("ERROR: [ログ取得API][JobID=%s]ログの読み込みに失敗しました", jobId)
			return nil, fmt.Errorf("レスポンスボディ取得に失敗しました: %w", err)
		}
		log.Printf("INFO: [ログ取得API][JobID=%s]ログの取得に成功しました", jobId)
		return &result, nil
		// エラー応答
	} else {
		result := dto.GitHubErrorResponse{}
		json.NewDecoder(response.Body).Decode(&result)
		log.Printf("ERROR: [ログ取得API][JobID=%s]ログの取得に失敗しました", jobId)
		log.Printf("ERROR: [ログ取得API][JobID=%s]エラーコード: %s", jobId, response.Status)
		log.Printf("ERROR: [ログ取得API][JobID=%s]レスポンスエラー: %s", jobId, result.Message)
		log.Printf("ERROR: [ログ取得API][JobID=%s]ドキュメント: %s", jobId, result.DocumentationUrl)
		return nil, fmt.Errorf("GitHub API Error : %s", result.Message)
	}
}

// OrganizationTokenを取得する
func getOrganizationToken() (string, error) {
	secretPath := "/run/secrets/"
	organizationTokenSecretName := "organization-token"
	organizationToken := ""
	// シークレット(/run/secrets/*)を探す
	fp, err := os.Open(secretPath + organizationTokenSecretName)
	if err == nil {
		log.Printf("INFO: Organizationトークンをシークレットから取得します")
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Scan()
		organizationToken = scanner.Text()
	} else if organizationToken == "" {
		// 見つからない場合は環境変数を探す
		log.Printf("INFO: Organizationトークンを環境変数から取得します")
		organizationToken = os.Getenv(strings.ToUpper(organizationTokenSecretName))
	}

	if organizationToken == "" {
		log.Printf("ERROR: Organizationトークンが見つかりませんでした")
		return organizationToken, fmt.Errorf("Organizationトークンが見つかりませんでした")
	}
	return organizationToken, nil
}

// dto.JobのJobIDフィールドをgitHub.Jobにコピーする
func convertGitHubJobListToJobList(githubJobList *[]dto.GitHubJob) *[]github.Job {
	jobList := []github.Job{}
	for _, githubJob := range *githubJobList {
		job := new(github.Job)
		log.Printf("jobId= %d", githubJob.JobId)
		job.JobId = strconv.Itoa(githubJob.JobId)
		job.RunId = strconv.Itoa(githubJob.RunId)
		job.JobName = githubJob.JobName
		job.Status = strings.ToUpper(githubJob.Status)
		jobList = append(jobList, *job)
	}
	return &jobList
}
