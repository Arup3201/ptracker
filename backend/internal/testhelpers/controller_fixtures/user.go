package controller_fixtures

type UserParams struct {
	IDPSubject  string
	IDPProvider string
	Username    string
	Email       string
	DisplayName *string
	AvatarURL   *string
}

func (f *ControllerFixtures) User(user UserParams) string {
	id, err := f.store.User().Create(f.ctx, user.IDPSubject, user.IDPProvider, user.Username,
		user.Email, user.DisplayName, user.AvatarURL)
	if err != nil {
		panic("store user create error")
	}

	return id
}
