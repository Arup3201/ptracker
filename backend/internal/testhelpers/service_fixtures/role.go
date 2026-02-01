package service_fixtures

import "fmt"

func (f *Fixtures) Role(projectID, userID, role string) {
	err := f.store.Role().Create(
		f.ctx,
		projectID,
		userID,
		role,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create role fixture: %v", err))
	}
}
