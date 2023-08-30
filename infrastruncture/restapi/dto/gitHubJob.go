package dto

type GitHubJob struct {
	JobId       int    `json:"id"`
	RunId       int    `json:"run_id"`
	JobName     string `json:"name"`
	Status      string `json:"status"`
	Conclustion string `json:"conclusion"`
}
