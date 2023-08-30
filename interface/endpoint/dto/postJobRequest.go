package dto

type PostJobRequest struct {
	RepositoryId   string `json:"repository_id"`
	RepositoryName string `json:"repository_name"`
	RunId          string `json:"run_id"`
	WorkflowRef    string `json:"workflow_ref"`
	JobName        string `json:"job_name"`
	RunAttempt     string `json:"run_attempt"`
	Status         string `json:"status"`
}
