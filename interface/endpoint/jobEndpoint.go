package endpoint

import (
	"net/http"

	"github.com/kazuo278/dashboard/application"
	"github.com/kazuo278/dashboard/interface/endpoint/dto"
	"github.com/kazuo278/dashboard/interface/websocket"
	"github.com/labstack/echo/v4"
)

type JobEndpoint interface {
	PostJob(c echo.Context) error
	PutJob(c echo.Context) error
}

type jobEndpointImpl struct {
	jobApp application.JobApplication
}

func NewJobEndpoint(jobApp application.JobApplication) JobEndpoint {
	jobEndpoint := new(jobEndpointImpl)
	jobEndpoint.jobApp = jobApp
	return jobEndpoint
}

// 実行履歴を登録する
// POST: /actions/history
// JSON:
//
//	{
//		run_id: <string>,
//		run_attempt: <string>,
//		repository_id: <string>,
//		repository_name: <string>,
//		workflow_ref: <string>,
//		job_name: <string>s
//	}
func (endpoint *jobEndpointImpl) PostJob(c echo.Context) error {
	// JSONリクエストを取得
	body := dto.PostJobRequest{}
	c.Bind(&body)
	// 実行履歴を登録
	err := endpoint.jobApp.SetUpRunner(body.RepositoryId, body.RepositoryName, body.RunId, body.WorkflowRef, body.JobName, body.RunAttempt)
	// 登録失敗エラーハンドリング
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	// WebSocketを確立したブラウザへ更新を通知
	websocket := websocket.NewDashboardWebSocket()
	// ブラウザ更新処理を非同期呼び出し
	go websocket.Update()

	return c.JSON(http.StatusNoContent, nil)
}

// 実行履歴を更新する
// PUT: /actions/history
// JSON:
//
//	{
//		repository_id: <string>,
//	  run_id: <string>,
//		job_name: <string>,
//		run_attempt: <string>
//	}
func (endpoint *jobEndpointImpl) PutJob(c echo.Context) error {
	// JSONリクエストを取得
	body := dto.PutJobRequest{}
	c.Bind(&body)
	// 実行履歴とジョブ解析を実行
	err := endpoint.jobApp.CompletedRunner(body.RepositoryId, body.RunId, body.JobName, body.RunAttempt)
	// 登録失敗エラーハンドリング
	if err != nil {
		return c.JSON(http.StatusInternalServerError, err)
	}
	// WebSocketを確立したブラウザへ更新を通知
	websocket := websocket.NewDashboardWebSocket()
	// ブラウザ更新処理を非同期呼び出し
	go websocket.Update()
	return c.JSON(http.StatusNoContent, nil)
}
