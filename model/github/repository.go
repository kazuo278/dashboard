package github

type Repository struct {
	RepositoryId   string `gorm:"primaryKey;column:repository_id" json:"repository_id"`
	RepositoryName string `gorm:"column:repository_name" json:"repository_name"`
}
