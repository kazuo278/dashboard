package dto

import "github.com/kazuo278/dashboard/application/dto"

type GetJobDetailsResponse struct {
	Count      int                  `json:"count"`
	JobDetails *[]dto.RepoJobDetail `json:"details"`
}
