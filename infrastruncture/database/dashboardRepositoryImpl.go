package database

import (
	"database/sql"
	"log"
	"os"

	"github.com/kazuo278/dashboard/application/dto"
	"github.com/kazuo278/dashboard/application/repository"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type dashboardRepositoryImpl struct {
	db *gorm.DB
}

func NewDashboardRepository() repository.DashboardRepository {
	dashboardRepository := new(dashboardRepositoryImpl)
	connectionUrl := os.Getenv("DATABASE_URL")
	sqlDB, err := sql.Open("pgx", connectionUrl)
	if err != nil {
		log.Print("ERROR: DB接続に失敗しました")
		panic(err)
	}
	dashboardRepository.db, err = gorm.Open(postgres.New(postgres.Config{
		Conn: sqlDB,
	}), &gorm.Config{})
	if err != nil {
		log.Print("ERROR: DB接続に失敗しました")
		panic(err)
	}
	dashboardRepository.db.Logger = dashboardRepository.db.Logger.LogMode(logger.Info)
	return dashboardRepository
}

// Job検索
func (repo dashboardRepositoryImpl) GetJobs(
	limit int, offset int, jobId string, repositoryId string, repositoryName string,
	workflowRef string, jobName string, runId string, runAttempt string, status string,
	startedAt string, finishedAt string) (*[]dto.RepoJob, int) {
	result := []dto.RepoJob{}
	count := 0
	sql := repo.db.Table("jobs").Select("jobs.*, repositories.repository_name").Joins("left join repositories on repositories.repository_id = jobs.repository_id")
	countSql := repo.db.Table("jobs").Select("count(1)").Joins("left join repositories on repositories.repository_id = jobs.repository_id")

	if jobId != "" {
		sql.Where("jobs.job_id = ?", repositoryId)
		countSql.Where("jobs.job_id = ?", repositoryId)
	}

	if repositoryId != "" {
		sql.Where("jobs.repository_id = ?", repositoryId)
		countSql.Where("jobs.repository_id = ?", repositoryId)
	}

	if repositoryName != "" {
		sql.Where("repository_name LIKE ?", repositoryName+"%")
		countSql.Where("repository_name LIKE ?", repositoryName+"%")
	}

	if workflowRef != "" {
		sql.Where("workflow_ref LIKE ?", workflowRef+"%")
		countSql.Where("workflow_ref LIKE ?", workflowRef+"%")
	}

	if jobName != "" {
		sql.Where("job_name LIKE ?", jobName+"%")
		countSql.Where("job_name LIKE ?", jobName+"%")
	}

	if runId != "" {
		sql.Where("jobs.run_id = ?", repositoryId)
		countSql.Where("jobs.run_id = ?", repositoryId)
	}

	if runAttempt != "" {
		sql.Where("run_attempt = ? ", runAttempt)
		countSql.Where("run_attempt = ?", runAttempt)
	}

	if status != "" {
		sql.Where("status = ?", status)
		countSql.Where("status = ?", status)
	}

	if startedAt != "" {
		sql.Where("started_at >= ?", startedAt)
		countSql.Where("started_at >= ?", startedAt)
	}

	if finishedAt != "" {
		sql.Where("finished_at <= ?", finishedAt)
		countSql.Where("finished_at <= ?", finishedAt)
	}

	sql.Order("jobs.started_at desc").Limit(limit).Offset(offset).Scan(&result)
	countSql.Scan(&count)

	return &result, count
}

// リポジトリごとのJob実行回数検索
func (repo dashboardRepositoryImpl) GetJobCount(repositoryName string, startedAt string, finishedAt string) (*[]dto.RepoCount, int) {
	result := []dto.RepoCount{}
	var totalCount int
	sql := repo.db.Table("jobs").Select("repositories.repository_name, count(1) as count")
	totalSql := repo.db.Table("jobs").Select("count(*) as total")
	sql.Joins("left join repositories on repositories.repository_id = jobs.repository_id")
	totalSql.Joins("left join repositories on repositories.repository_id = jobs.repository_id")

	if repositoryName != "" {
		sql.Where("repository_name LIKE ?", repositoryName+"%")
		totalSql.Where("repository_name LIKE ?", repositoryName+"%")
	}

	if startedAt != "" {
		sql.Where("started_at >= ?", startedAt)
		totalSql.Where("started_at >= ?", startedAt)
	}

	if finishedAt != "" {
		sql.Where("finished_at <= ?", finishedAt)
		totalSql.Where("finished_at <= ?", finishedAt)
	}

	sql.Group("repositories.repository_name").Order("count desc").Scan(&result)
	totalSql.Scan(&totalCount)

	return &result, totalCount
}

// リポジトリごとのJob実行時間検索
func (repo dashboardRepositoryImpl) GetJobTime(repositoryName string, startedAt string, finishedAt string) (*[]dto.RepoTime, int) {
	result := []dto.RepoTime{}
	var totalSeconds int

	sql := repo.db.Table("jobs").Select("repositories.repository_name, round(sum(extract(epoch from finished_at) - extract(epoch from started_at))) as seconds")
	totalSql := repo.db.Table("jobs").Select("coalesce(round(sum(extract(epoch from finished_at) - extract(epoch from started_at))), 0) as total")

	sql.Joins("left join repositories on repositories.repository_id = jobs.repository_id")
	totalSql.Joins("left join repositories on repositories.repository_id = jobs.repository_id")

	sql.Where("status = ?", "COMPLETED")
	totalSql.Where("status = ?", "COMPLETED")

	if repositoryName != "" {
		sql.Where("repository_name LIKE ?", repositoryName+"%")
		totalSql.Where("repository_name LIKE ?", repositoryName+"%")
	}

	if startedAt != "" {
		sql.Where("started_at >= ?", startedAt)
		totalSql.Where("started_at >= ?", startedAt)
	}

	if finishedAt != "" {
		sql.Where("finished_at <= ?", finishedAt)
		totalSql.Where("finished_at <= ?", finishedAt)
	}

	sql.Group("repositories.repository_name").Order("seconds desc").Scan(&result)
	totalSql.Scan(&totalSeconds)

	return &result, totalSeconds
}

// Job検索
func (repo dashboardRepositoryImpl) GetJobDetails(
	limit int, offset int, jobId string, repositoryId string, repositoryName string,
	usingPath string, usingRef string, jobName string, runId string, runAttempt string, typeName string, startedAt string,
	finishedAt string) (*[]dto.RepoJobDetail, int) {
	result := []dto.RepoJobDetail{}
	count := 0
	sql := repo.db.Table("job_details").Select("job_details.job_detail_id, job_details.job_id, job_details.using_path, job_details.using_ref, job_details.type, jobs.run_id, jobs.run_attempt, jobs.job_name, jobs.started_at, jobs.finished_at, repositories.repository_id, repositories.repository_name")
	sql.Joins("left join jobs on jobs.job_id = job_details.job_id")
	sql.Joins("left join repositories on repositories.repository_id = jobs.repository_id")
	countSql := repo.db.Table("job_details").Select("count(*) as total")
	countSql.Joins("left join jobs on jobs.job_id = job_details.job_id")
	countSql.Joins("left join repositories on repositories.repository_id = jobs.repository_id")

	if jobId != "" {
		sql.Where("job_details.job_id = ?", jobId)
		countSql.Where("job_details.job_id = ?", jobId)
	}

	if repositoryId != "" {
		sql.Where("repositories.repository_id = ?", repositoryId)
		countSql.Where("repositories.repository_id = ?", repositoryId)
	}

	if repositoryName != "" {
		sql.Where("repository_name LIKE ?", repositoryName+"%")
		countSql.Where("repository_name LIKE ?", repositoryName+"%")
	}

	if usingPath != "" {
		sql.Where("using_path LIKE ?", usingPath+"%")
		countSql.Where("using_path LIKE ?", usingPath+"%")
	}

	if usingRef != "" {
		sql.Where("using_ref LIKE ?", usingRef+"%")
		countSql.Where("using_ref LIKE ?", usingRef+"%")
	}

	if jobName != "" {
		sql.Where("job_name LIKE ?", jobName+"%")
		countSql.Where("job_name LIKE ?", jobName+"%")
	}

	if runId != "" {
		sql.Where("run_id = ? ", runId)
		countSql.Where("run_id = ?", runId)
	}

	if runAttempt != "" {
		sql.Where("run_attempt = ? ", runAttempt)
		countSql.Where("run_attempt = ?", runAttempt)
	}

	if typeName != "" {
		sql.Where("type = ? ", typeName)
		countSql.Where("type = ?", typeName)
	}

	if startedAt != "" {
		sql.Where("started_at >= ?", startedAt)
		countSql.Where("started_at >= ?", startedAt)
	}

	if finishedAt != "" {
		sql.Where("finished_at <= ?", finishedAt)
		countSql.Where("finished_at <= ?", finishedAt)
	}

	sql.Order("started_at desc").Limit(limit).Offset(offset).Scan(&result)
	countSql.Scan(&count)

	return &result, count
}
