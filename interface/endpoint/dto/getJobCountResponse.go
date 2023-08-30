package dto

import "github.com/kazuo278/dashboard/application/dto"

type GetJobCountResponse struct {
	Repositories *[]dto.RepoCount `json:"repositories"`
	TotalCount   int              `json:"total_count"`
}
