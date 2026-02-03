package repo_fixtures

type CommentParam struct {
	ProjectId string
	TaskId    string
	UserId    string
	Content   string
}

func GetCommentRow(projectId, taskId, userId, content string) CommentParam {
	return CommentParam{
		ProjectId: projectId,
		TaskId:    taskId,
		UserId:    userId,
		Content:   content,
	}
}

func (f *Fixtures) InsertComment(c CommentParam) {
	_, err := f.db.ExecContext(
		f.ctx,
		`
		INSERT INTO comments (project_id, task_id, user_id, content)
		VALUES ($1,$2,$3,$4)
		`,
		c.ProjectId,
		c.TaskId,
		c.UserId,
		c.Content,
	)
	if err != nil {
		panic("insert comment fixture failed: " + err.Error())
	}
}
