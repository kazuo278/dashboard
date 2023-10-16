package main

import (
	"github.com/kazuo278/dashboard/application"
	"github.com/kazuo278/dashboard/infrastruncture/database"
	"github.com/kazuo278/dashboard/infrastruncture/restapi"
	"github.com/kazuo278/dashboard/interface/endpoint"
	"github.com/kazuo278/dashboard/interface/websocket"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func main() {
	// インフラ層の初期化
	jobRepository := database.NewJobRepository()
	dashboardRepository := database.NewDashboardRepository()
	jobApi := restapi.NewGitHubJobApi()
	// アプリケーション層の初期化
	jobApp := application.NewJobAppication(jobRepository, jobApi)
	dashboardApp := application.NewDashboardApplication(dashboardRepository)
	application.NewReconcileAppication(jobRepository, jobApi)

	// インターフェース層の初期化
	e := echo.New()
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	// REST API
	// 更新系API
	jobEndpoint := endpoint.NewJobEndpoint(jobApp)
	e.POST("/actions/history", jobEndpoint.PostJob)
	e.PUT("/actions/history", jobEndpoint.PutJob)
	// 参照系API
	dashboardEndpoint := endpoint.NewDashboardEndpoint(dashboardApp)
	e.GET("/actions", dashboardEndpoint.GetJobs)
	e.GET("/actions/count", dashboardEndpoint.GetJobCount)
	e.GET("/actions/time", dashboardEndpoint.GetJobTime)
	e.GET("/details", dashboardEndpoint.GetJobDetails)
	// WebContent
	dashboardWebSocket := websocket.NewDashboardWebSocket()
	e.GET("/ws", dashboardWebSocket.Socket)
	e.Static("/dashboard", "./static-content")
	// WebServer起動
	e.Logger.Fatal(e.Start(":8080"))
}
