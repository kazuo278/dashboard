package database

import (
	"log"
	"os"
	"time"

	"database/sql"

	"github.com/kazuo278/dashboard/domain/model/github"
	"github.com/kazuo278/dashboard/domain/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type jobRepositoryImpl struct {
	db *gorm.DB
}

func NewJobRepository() repository.JobRepository {
	jobRepository := new(jobRepositoryImpl)
	connectionUrl := os.Getenv("DATABASE_URL")
	sqlDB, err := sql.Open("pgx", connectionUrl)
	if err != nil {
		log.Print("ERROR: DB接続に失敗しました")
		panic(err)
	}
	jobRepository.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Print("ERROR: DB接続に失敗しました")
		panic(err)
	}
	jobRepository.db.Logger = jobRepository.db.Logger.LogMode(logger.Info)
	return jobRepository
}

// RepositoryのID検索
func (repo *jobRepositoryImpl) GetRepositoryById(repositoryId string) *github.Repository {
	result := github.Repository{}
	repo.db.Where("repository_id = ?", repositoryId).Limit(1).Find(&result)
	return &result
}

// リポジトリの登録
func (repo *jobRepositoryImpl) CreateRepository(repository *github.Repository) *github.Repository {
	repo.db.Save(repository)
	return repository
}

// ID検索
func (repo *jobRepositoryImpl) GetJobById(jobId string) *github.Job {
	result := github.Job{}
	repo.db.Where("job_id = ?", jobId).Limit(1).Find(&result)
	return &result
}

// ジョブの登録
func (repo *jobRepositoryImpl) CreateJob(job *github.Job) *github.Job {
	repo.db.Create(job)
	return job
}

// ジョブの更新
func (repo *jobRepositoryImpl) UpdateJob(job *github.Job) *github.Job {
	repo.db.Save(job)
	return job
}

// ジョブ詳細の登録
func (repo *jobRepositoryImpl) CreateJobDetail(detail *github.JobDetail) *github.JobDetail {
	repo.db.Save(detail)
	return detail
}

// 未終了ジョブの取得
func (repo *jobRepositoryImpl) GetUnfinishedJob(from time.Time, to time.Time) *github.Job {
	result := github.Job{}
	sql := repo.db.Where("status != ?", github.STATUS_COMPLETED)
	sql.Where("started_at <= ?", to)
	sql.Where("started_at >= ?", from).Limit(1).Find(&result)
	return &result
}
