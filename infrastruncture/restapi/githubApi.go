package restapi

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/kazuo278/action-dashboard/domain/model/github"
	"github.com/kazuo278/action-dashboard/infrastruncture/restapi/dto"
)

const baseUrl string = "https://api.github.com/"
const versionHeader string = "X-GitHub-Api-Version"
const version string = "2022-11-28"
const AcceptHeader string = "Accept"
const accept string = "application/vnd.github+json"
const AuthorizationHeader string = "Authorization"

type GitHubJobAPI struct {
}

// ジョブIDリストの取得
func (api GitHubJobAPI) GetJobList(runId string, repositoryName string, runAttempt string) (*[]github.Job, error) {
	// 参考：https://docs.github.com/ja/rest/actions/workflow-jobs?apiVersion=2022-11-28#list-jobs-for-a-workflow-run-attempt
	// URI：https://api.github.com/repos/OWNER/REPO/actions/runs/RUN_ID/attempts/ATTEMPT_NUMBER/jobs
	jobListEndpointUri := baseUrl + repositoryName + "/actions/runs/" + runId + "/attempts/" + runAttempt + "/jobs"

	// トークン取得
	organizationToken, err := getOrganizationToken()
	if err != nil {
		return nil, err
	}

	// リクエスト送信
	request, _ := http.NewRequest(http.MethodPost, jobListEndpointUri, nil)
	request.Header.Set(versionHeader, version)
	request.Header.Set(AcceptHeader, accept)
	request.Header.Set(AuthorizationHeader, "Bearer "+organizationToken)
	client := new(http.Client)

	// レスポンス受信
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error: エンドポイント接続時にエラーが発生しました。%s", err)
		return nil, fmt.Errorf("error: エンドポイント接続時にエラーが発生しました: %w", err)
	}

	defer response.Body.Close()
	// 正常応答
	if response.StatusCode == 200 {
		// レスポンスのJSON文字列を構造体GitHubResponseにマッピング
		result := dto.GitHubJobList{}
		json.NewDecoder(response.Body).Decode(&result)
		return convertGitHubJobListToJobList(result.Jobs), nil
		// エラー応答
	} else {
		log.Printf("レスポンスエラー: %s", err)
		result := dto.GitHubErrorResponse{}
		json.NewDecoder(response.Body).Decode(&result)
		return nil, fmt.Errorf("GitHub API Error : %s", result.Message)
	}
}

// ジョブリストの取得
func GetJob(jobId string, repositoryName string) (*github.Job, error) {
	// 参考：https://docs.github.com/ja/rest/actions/workflow-jobs?apiVersion=2022-11-28#get-a-job-for-a-workflow-run
	// URI：https://api.github.com/repos/OWNER/REPO/actions/jobs/JOB_ID
	jobListEndpointUri := baseUrl + repositoryName + "/actions/jobs/" + jobId

	// トークン取得
	organizationToken, err := getOrganizationToken()
	if err != nil {
		return nil, err
	}

	// リクエスト送信
	request, _ := http.NewRequest(http.MethodPost, jobListEndpointUri, nil)
	request.Header.Set(versionHeader, version)
	request.Header.Set(AcceptHeader, accept)
	request.Header.Set(AuthorizationHeader, "Bearer "+organizationToken)
	client := new(http.Client)

	// レスポンス受信
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error: エンドポイント接続時にエラーが発生しました。%s", err)
		return nil, fmt.Errorf("error: エンドポイント接続時にエラーが発生しました: %w", err)
	}

	defer response.Body.Close()
	// 正常応答
	if response.StatusCode == 200 {
		// レスポンスのJSON文字列を構造体GitHubResponseにマッピング
		result := github.Job{}
		json.NewDecoder(response.Body).Decode(&result)
		return &result, nil
		// エラー応答
	} else {
		log.Printf("レスポンスエラー: %s", err)
		result := dto.GitHubErrorResponse{}
		json.NewDecoder(response.Body).Decode(&result)
		return nil, fmt.Errorf("GitHub API Error : %s", result.Message)
	}
}

// JobIdに紐づくログを取得する
func GetJobLog(jobId string, repositoryName string) (*[]byte, error) {
	// 参考：https://docs.github.com/ja/rest/actions/workflow-jobs?apiVersion=2022-11-28#download-job-logs-for-a-workflow-run
	// URI：https://api.github.com/repos/OWNER/REPO/actions/jobs/JOB_ID/logs
	jobListEndpointUri := baseUrl + repositoryName + "/actions/jobs/" + jobId + "/logs"

	// トークン取得
	organizationToken, err := getOrganizationToken()
	if err != nil {
		return nil, err
	}

	// リクエスト送信
	request, _ := http.NewRequest(http.MethodPost, jobListEndpointUri, nil)
	request.Header.Set(versionHeader, version)
	request.Header.Set(AcceptHeader, accept)
	request.Header.Set(AuthorizationHeader, "Bearer "+organizationToken)
	client := new(http.Client)

	// レスポンス受信
	// 本来302レスポンスを受けリダイレクトする必要があるが、net/httpの仕組みにより自動でリダイレクトする。
	response, err := client.Do(request)
	if err != nil {
		log.Printf("Error: エンドポイント接続時にエラーが発生しました。%s", err)
		return nil, fmt.Errorf("error: エンドポイント接続時にエラーが発生しました: %w", err)
	}

	defer response.Body.Close()
	// 正常応答
	if response.StatusCode == 200 {
		// レスポンスのログ文字列を返却
		result, err := io.ReadAll(response.Body)
		if err != nil {
			log.Printf("レスポンスボディ取得に失敗しました: %s", err)
			return nil, fmt.Errorf("レスポンスボディ取得に失敗しました: %w", err)
		}
		return &result, nil
		// エラー応答
	} else {
		log.Printf("レスポンスエラー: %s", err)
		result := dto.GitHubErrorResponse{}
		json.NewDecoder(response.Body).Decode(&result)
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
		log.Printf("Organizationトークンをシークレットから取得します")
		defer fp.Close()
		scanner := bufio.NewScanner(fp)
		scanner.Scan()
		organizationToken = scanner.Text()
	} else if organizationToken == "" {
		// 見つからない場合は環境変数を探す
		log.Printf("Organizationトークンを環境変数から取得します")
		organizationToken = os.Getenv(strings.ToUpper(organizationTokenSecretName))
	}

	if organizationToken == "" {
		return organizationToken, fmt.Errorf("Organizationトークンが見つかりませんでした")
	}
	return organizationToken, nil
}

// dto.JobのJobIDフィールドをgitHub.Jobにコピーする
func convertGitHubJobListToJobList(githubJobList *[]dto.GitHubJob) *[]github.Job {
	jobList := []github.Job{}
	for _, githubJob := range *githubJobList {
		job := new(github.Job)
		job.JobId = githubJob.JobId
		jobList = append(jobList, *job)
	}
	return &jobList
}
