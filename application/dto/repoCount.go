package dto

type RepoCount struct {
	RepositoryName string `gorm:"column:repository_name" json:"repository_name"`
	Count          int    `gorm:"column:count" json:"count"`
}
