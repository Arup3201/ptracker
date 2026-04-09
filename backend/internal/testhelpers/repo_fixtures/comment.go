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
	if f.db != nil {
		if err := f.db.WithContext(f.ctx).Create(&c).Error; err != nil {
			panic("insert comment fixture failed: " + err.Error())
		}
		return
	}
}
