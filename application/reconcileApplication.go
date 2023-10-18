package application

import (
	"log"
	"time"

	"github.com/kazuo278/dashboard/domain/model/github"
	"github.com/kazuo278/dashboard/domain/repository"
	"github.com/kazuo278/dashboard/domain/service"
)

type ReconcileApplication interface {
	ReconcileAllUnfinishedJobsWithin24h() error
	ReconcileSingleUnfinishedJob(job *github.Job, repositoryId string) error
}

type reconcileApplicationImpl struct {
	analyzer      analyzer
	jobRepository repository.JobRepository
	jobApi        service.JobApi
}

func NewReconcileAppication(jobRepositroy repository.JobRepository, jobApi service.JobApi) ReconcileApplication {
	reconcileApplication := new(reconcileApplicationImpl)
	reconcileApplication.jobRepository = jobRepositroy
	reconcileApplication.jobApi = jobApi
	reconcileApplication.analyzer = newAnalyzer(jobRepositroy, jobApi)
	go func() {
		for {
			time.Sleep(60 * time.Second)
			reconcileApplication.ReconcileAllUnfinishedJobsWithin24h()
		}
	}()
	return reconcileApplication
}

// 1日以上未終了なジョブの情報を更新する
func (app *reconcileApplicationImpl) ReconcileAllUnfinishedJobsWithin24h() error {
	// 1日以上未終了なジョブを取得
	toDate := time.Now().Add(time.Hour * (-24))
	// 90日以上経過したジョブは対象外
	fromDate := time.Now().Add(time.Hour * (-24 * 90))
	// 90日前〜1日前に開始され、未完了のジョブを取得
	unfinishedJob := app.jobRepository.GetUnfinishedJob(fromDate, toDate)
	if (*unfinishedJob == github.Job{}) {
		// 未終了のジョブが存在しない場合は操作なし
		return nil
	}
	// 整合チェック
	app.ReconcileSingleUnfinishedJob(unfinishedJob, unfinishedJob.RepositoryId)
	return nil
}

func (app *reconcileApplicationImpl) ReconcileSingleUnfinishedJob(oldJob *github.Job, repositoryId string) error {
	log.Printf("INFO: [整合チェック][Starting][JobID=%s]ジョブを更新します", oldJob.JobId)
	repository := app.jobRepository.GetRepositoryById(repositoryId)
	unfinishedJob, err := app.jobApi.GetJob(oldJob.JobId, repository.RepositoryName)
	if err != nil {
		return err
	}
	// APIでは連携されない情報をセット
	unfinishedJob.RunAttempt = oldJob.RunAttempt
	unfinishedJob.RepositoryId = repositoryId
	unfinishedJob.WorkflowRef = oldJob.WorkflowRef
	// APIから取得した情報で更新
	app.jobRepository.UpdateJob(unfinishedJob)
	app.analyzer.analyze(unfinishedJob.JobId, repository.RepositoryName)
	log.Printf("INFO: [整合チェック][Finished][JobID=%s]ジョブを更新完了しました", unfinishedJob.JobId)
	return nil
}
