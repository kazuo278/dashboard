package database

import (
	"log"
	"os"

	"database/sql"

	"github.com/kazuo278/action-dashboard/domain/model/github"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type JobRepositoryImpl struct {
	db *gorm.DB
}

func Init() *JobRepositoryImpl {
	jobRepository := new(JobRepositoryImpl)
	connectionUrl := os.Getenv("DATABASE_URL")
	var err error
	sqlDB, err := sql.Open("pgx", connectionUrl)
	if err != nil {
		log.Print("DB接続に失敗しました")
		panic(err)
	}
	jobRepository.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Print("DB接続に失敗しました")
		panic(err)
	}
	return jobRepository
}

// RepositoryのID検索
func (repo JobRepositoryImpl) GetRepositoryById(repositoryId string) *github.Repository {
	result := github.Repository{}
	repo.db.Where("repository_id = ?", repositoryId).Limit(1).Find(&result)
	return &result
}

// リポジトリの登録
func (repo JobRepositoryImpl) CreateRepository(repository *github.Repository) *github.Repository {
	repo.db.Save(repository)
	return repository
}

// ID代替の複数条件によるジョブの１件取得
func (repo JobRepositoryImpl) GetJobByIds(repositoryId string, runId string, jobName string, runAttempt string) *github.Job {
	result := github.Job{}
	repo.db.Where("repository_id = ? AND run_id = ? AND job_name = ? AND run_attempt = ?", repositoryId, runId, jobName, runAttempt).First(&result)
	return &result
}

// ジョブの登録
func (repo JobRepositoryImpl) CreateJob(job *github.Job) *github.Job {
	repo.db.Create(job)
	return job
}

// ジョブの更新
func (repo JobRepositoryImpl) UpdateJob(job *github.Job) *github.Job {
	repo.db.Save(job)
	return job
}
