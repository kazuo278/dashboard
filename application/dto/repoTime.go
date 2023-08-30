package dto

type RepoTime struct {
	RepositoryName string `gorm:"column:repository_name" json:"repository_name"`
	Seconds        int    `gorm:"column:seconds" json:"seconds"`
}
