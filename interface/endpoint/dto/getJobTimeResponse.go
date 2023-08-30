package dto

import "github.com/kazuo278/dashboard/application/dto"

type GetJobTimeResponse struct {
	Repositories *[]dto.RepoTime `json:"repositories"`
	TotalSeconds int             `json:"total_seconds"`
}
