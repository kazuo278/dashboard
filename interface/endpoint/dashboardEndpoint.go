package endpoint

import (
	"net/http"
	"strconv"

	"github.com/kazuo278/dashboard/application"
	"github.com/kazuo278/dashboard/interface/endpoint/dto"
	"github.com/labstack/echo/v4"
)

type DashboardEndpoint interface {
	GetJobs(c echo.Context) error
	GetJobCount(c echo.Context) error
	GetJobTime(c echo.Context) error
	GetJobDetails(c echo.Context) error
}

type dashboardEndpointImpl struct {
	dashboardApp application.DashboardApplication
}

func NewDashboardEndpoint(app application.DashboardApplication) DashboardEndpoint {
	endpoint := new(dashboardEndpointImpl)
	endpoint.dashboardApp = app
	return endpoint
}

// 指定された件数の実行履歴を返す
// GET: /actions
// QueryParam:
//
//	limit=<int>
//	offset=<int>
//	job_id=<string>
//	repository_id=<string>
//	repository_name=<string>
//	workflow_ref=<string>
//	job_name=<string>
//	run_id=<string>
//	run_attempt=<string>
//	status=<string>
//	conclusion=<string>
//	started_at=<string>
//	finished_at=<string>
func (endpoint *dashboardEndpointImpl) GetJobs(c echo.Context) error {
	limitNum, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limitNum >= 100 {
		limitNum = 50
	}
	offsetNum, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offsetNum == 0 {
		offsetNum = 0
	}
	jobs, count := endpoint.dashboardApp.GetJobs(limitNum, offsetNum, c.FormValue("job_id"), c.QueryParam("repository_id"), c.QueryParam("repository_name"), c.QueryParam("workflow_ref"), c.QueryParam("job_name"), c.QueryParam("run_id"), c.QueryParam("run_attempt"), c.QueryParam("status"), c.QueryParam("conclusion"), c.QueryParam("started_at"), c.QueryParam("finished_at"))

	result := new(dto.GetJobsResponse)
	result.Jobs = jobs
	result.Count = count
	return c.JSON(http.StatusOK, result)
}

// リポジトリごとの実行回数を取得する
// GET: /actions/count
// QueryParam:
//
//	repository_name=<string>
//	started_at=<string>
//	finished_at=<string>
func (endpoint *dashboardEndpointImpl) GetJobCount(c echo.Context) error {
	// 実行回数を取得
	repos, count := endpoint.dashboardApp.GetJobCount(c.QueryParam("repository_name"), c.QueryParam("started_at"), c.QueryParam("finished_at"))
	result := new(dto.GetJobCountResponse)
	result.Repositories = repos
	result.TotalCount = count
	return c.JSON(http.StatusOK, result)
}

// リポジトリごとの実行時間(秒)を取得する
// GET: /actions/time
// QueryParam:
//
//	repository_name=<string>
//	started_at=<string>
//	finished_at=<string>
func (endpoint *dashboardEndpointImpl) GetJobTime(c echo.Context) error {
	// 実行回数を取得
	repos, seconds := endpoint.dashboardApp.GetJobTime(c.QueryParam("repository_name"), c.QueryParam("started_at"), c.QueryParam("finished_at"))
	result := new(dto.GetJobTimeResponse)
	result.Repositories = repos
	result.TotalSeconds = seconds
	return c.JSON(http.StatusOK, result)
}

// 指定された件数の実行詳細を返す
// GET: /details
// QueryParam:
//
//	limit=<int>
//	offset=<int>
//	job_id=<string>
//	repository_id=<string>
//	repository_name=<string>
//	using_path=<string>
//	using_ref=<string>
//	job_name=<string>
//	run_id=<string>
//	run_attempt=<string>
//	type=<string>
//	started_at=<string>
//	finished_at=<string>
func (endpoint *dashboardEndpointImpl) GetJobDetails(c echo.Context) error {
	limitNum, err := strconv.Atoi(c.QueryParam("limit"))
	if err != nil || limitNum >= 100 {
		limitNum = 50
	}
	offsetNum, err := strconv.Atoi(c.QueryParam("offset"))
	if err != nil || offsetNum == 0 {
		offsetNum = 0
	}
	details, count := endpoint.dashboardApp.GetJobDetails(limitNum, offsetNum, c.QueryParam("job_id"),
		c.QueryParam("repository_id"), c.QueryParam("repository_name"), c.QueryParam("using_path"), c.QueryParam("using_ref"),
		c.QueryParam("job_name"), c.QueryParam("run_id"), c.QueryParam("run_attempt"), c.QueryParam("type"),
		c.QueryParam("started_at"), c.QueryParam("finished_at"))
	result := new(dto.GetJobDetailsResponse)
	result.JobDetails = details
	result.Count = count
	return c.JSON(http.StatusOK, result)
}
