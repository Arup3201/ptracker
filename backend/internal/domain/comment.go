package domain

type Comment struct {
	Id        string  `json:"id"`
	ProjectId string  `json:"project_id"`
	TaskId    string  `json:"task_id"`
	User      *Member `json:"user"`
	Content   string  `json:"content"`
	CreatedAt int64   `json:"created_at"`
	UpdatedAt int64   `json:"updated_at"`
}
