package dto

type GitHubJob struct {
	JobId       string `json:"id"`
	RunId       string `json:"run_id"`
	Status      string `json:"status"`
	Conclustion string `json:"conclustion"`
}
