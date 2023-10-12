package github

type JobDetail struct {
	JobDetailId uint   `gorm:"primaryKey;column:job_detail_id" json:"job_detail_id"`
	JobId       string `gorm:"column:job_id" json:"job_id"`
	Type        string `gorm:"column:type" json:"type"`
	UsingPath   string `gorm:"column:using_path" json:"using_path"`
	UsingRef    string `gorm:"column:using_ref" json:"using_ref"`
}

const (
	TYPE_REUSABLE_WORKFLOW = "REUSABLE_WORKFLOW"
	TYPE_ACTION            = "ACTION"
)
