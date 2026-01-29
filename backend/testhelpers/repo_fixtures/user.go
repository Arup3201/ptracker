package repo_fixtures

import (
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/ptracker/domain"
)

func RandomUserRow() domain.User {
	uId := uuid.NewString()

	return domain.User{
		Id:          uId,
		IDPSubject:  uId,
		IDPProvider: "google",
		Username:    "user" + uId,
		Email:       "user@test.com",
		IsActive:    true,
	}
}

func (f *Fixtures) InsertUser(u domain.User) string {
	now := time.Now()
	_, err := f.db.ExecContext(
		f.ctx,
		`
		INSERT INTO users (
			id, idp_subject, idp_provider,
			username, email, display_name, avatar_url, 
			created_at, updated_at
		) VALUES ($1,$2,$3,$4,$5,$6,$7,$8,$8)
		`,
		u.Id,
		u.IDPSubject,
		u.IDPProvider,
		u.Username,
		u.Email,
		u.DisplayName,
		u.AvatarURL,
		now,
	)
	if err != nil {
		panic(fmt.Sprintf("insert user fixture failed: %v", err))
	}

	return u.Id
}
