package fixtures

import "fmt"

type UserParams struct {
	IDPSubject  string
	IDPProvider string
	Username    string
	Email       string
	DisplayName *string
	AvatarURL   *string
}

func DefaultUser() UserParams {
	return UserParams{
		IDPSubject:  "sub-123",
		IDPProvider: "google",
		Username:    "john_doe",
		Email:       "john@example.com",
	}
}

func (f *Fixtures) User(p UserParams) string {
	id, err := f.store.User().Create(
		f.ctx,
		p.IDPSubject,
		p.IDPProvider,
		p.Username,
		p.Email,
		p.DisplayName,
		p.AvatarURL,
	)
	if err != nil {
		panic(fmt.Sprintf("failed to create user fixture: %v", err))
	}

	return id
}
