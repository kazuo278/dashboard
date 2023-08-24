package endpoint

// 実行履歴を登録する
// POST: /actions/history
// JSON:
//
//	{
//	   job_id: <string>,
//	   run_id: <string>,
//	   run_attempt: <string>,
//	   repository_id: <string>,
//		 repository_name: <string>,
//	   workflow_ref: <string>,
//	   job_name: <string>,
//	   status: <string>
//	}
// func PostHistory(c echo.Context) error {
// 	// TODO: WebSocketを確立したブラウザへ更新を通知
// 	// JSONリクエストを取得
// 	body := github.Repository{}
// 	c.Bind(&body)
// 	// 実行履歴を登録
// 	result := application.CreateHistoryWithStarted(body.RepositoryId, body.RepositoryName, body.RunId, body.WorkflowRef, body.JobName, body.RunAttempt)
// 	return c.JSON(http.StatusCreated, result)
// }

// 指定された件数の実行履歴を返す
// GET: /actions/history
// QueryParam:
//
//	limit=<int>
//	offset=<int>
//	repository_id=<string>
//	repository_name=<string>
//	workflow_ref=<string>
//	job_name=<string>
//	run_attempt=<string>
//	status=<string>
//	started_at=<string>
//	finished_at=<string>
// func GetHistory(c echo.Context) error {
// 	limitNum, err := strconv.Atoi(c.QueryParam("limit"))
// 	if err != nil || limitNum >= 100 {
// 		limitNum = 50
// 	}
// 	offsetNum, err := strconv.Atoi(c.QueryParam("offset"))
// 	if err != nil || offsetNum == 0 {
// 		offsetNum = 0
// 	}
// 	result := application.GetHistories(limitNum, offsetNum, c.QueryParam("repository_id"), c.QueryParam("repository_name"), c.QueryParam("workflow_ref"), c.QueryParam("job_name"), c.QueryParam("run_attempt"), c.QueryParam("status"), c.QueryParam("started_at"), c.QueryParam("finished_at"))
// 	return c.JSON(http.StatusOK, result)
// }

// 実行履歴を更新する
// PUT: /actions/history
// { repository_id: <string>, run_id: <string>, job_name: <string>, run_attempt: <string> }
// func PutHistory(c echo.Context) error {
// 	// WebSocketを確立したブラウザへ更新を通知
// 	websocket.IsUpdated = true
// 	// JSONリクエストを取得
// 	body := custom.HistoryRepository{}
// 	c.Bind(&body)
// 	// 実行履歴を登録
// 	result := application.UpdateHistoryWithFinished(body.RepositoryId, body.RunId, body.JobName, body.RunAttempt)
// 	return c.JSON(http.StatusOK, result)
// }

// リポジトリごとの実行回数を取得する
// GET: /actions/count?repository_name=<string>&started_at=<string>&finished_at=<string>
// func GetHistoryCount(c echo.Context) error {
// 	// 実行回数を取得
// 	result := application.GetHistoryCount(c.QueryParam("repository_name"), c.QueryParam("started_at"), c.QueryParam("finished_at"))
// 	return c.JSON(http.StatusOK, result)
// }

// リポジトリごとの実行時間(秒)を取得する
// GET: /actions/time?repository_name=<string>&started_at=<string>&finished_at=<string>
// func GetHistoryTime(c echo.Context) error {
// 	// 実行回数を取得
// 	result := application.GetHistoryTime(c.QueryParam("repository_name"), c.QueryParam("started_at"), c.QueryParam("finished_at"))
// 	return c.JSON(http.StatusOK, result)
// }
