package service_fixtures

import "fmt"

func (f *Fixtures) Membership(projectID, userID, role string) {
	err := f.store.Membership().Create(
		f.ctx,
		projectID,
		userID,
		role,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create role fixture: %v", err))
	}
}
