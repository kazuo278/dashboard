package dto

type PutJobRequest struct {
	RepositoryId string `json:"repository_id"`
	RunId        string `json:"run_id"`
	JobName      string `json:"job_name"`
	RunAttempt   string `json:"run_attempt"`
}
