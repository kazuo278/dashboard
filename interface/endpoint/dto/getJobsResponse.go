package dto

import "github.com/kazuo278/dashboard/application/dto"

type GetJobsResponse struct {
	Count int            `json:"count"`
	Jobs  *[]dto.RepoJob `json:"jobs"`
}
