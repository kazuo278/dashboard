package github

type JobDetail struct {
	JobDetailId         uint   `gorm:"primaryKey;column:job_detaild_id" json:"job_detaild_id"`
	JobId               string `gorm:"column:job_id" json:"job_id"`
	Type                string `gorm:"column:type" json:"type"`
	UsingRepositoryName string `gorm:"column:using_repository_name" json:"using_repository_name"`
	UsingRepositoryRef  string `gorm:"column:using_repository_ref" json:"using_repository_ref"`
}
