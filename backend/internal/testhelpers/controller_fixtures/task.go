package controller_fixtures

type TaskParams struct {
	ProjectID   string
	Name        string
	Description string
	Status      string
}

func (f *ControllerFixtures) Task(task TaskParams) string {
	id, err := f.store.Task().Create(f.ctx,
		task.ProjectID,
		task.Name,
		task.Description,
		task.Status)
	if err != nil {
		panic("store task create error")
	}

	return id
}
