package database

import (
	"log"
	"os"

	"database/sql"

	"github.com/kazuo278/action-dashboard/model/github"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func Init() {
	connectionUrl := os.Getenv("DATABASE_URL")
	var err error
	sqlDB, err := sql.Open("pgx", connectionUrl)
	if err != nil {
		log.Print("DB接続に失敗しました")
		panic(err)
	}
	db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Print("DB接続に失敗しました")
		panic(err)
	}
}

// RepositoryのID検索
func GetRepositoryById(repositoryId string) *github.Repository {
	result := github.Repository{}
	db.Where("repository_id = ?", repositoryId).Limit(1).Find(&result)
	return &result
}

// リポジトリの登録
func CreateRepository(repository *github.Repository) *github.Repository {
	db.Save(repository)
	return repository
}

// ID代替の複数条件によるジョブの１件取得
func GetJobByIds(repositoryId string, runId string, jobName string, runAttempt string) *github.Job {
	result := github.Job{}
	db.Where("repository_id = ? AND run_id = ? AND job_name = ? AND run_attempt = ?", repositoryId, runId, jobName, runAttempt).First(&result)
	return &result
}

// ジョブの登録
func CreateJob(job *github.Job) *github.Job {
	db.Create(job)
	return job
}

// ジョブの更新
func UpdateJob(job *github.Job) *github.Job {
	db.Save(job)
	return job
}
