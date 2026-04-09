package repo_fixtures

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ptracker/internal/repositories/models"
)

func RandomUserRow() models.User {
	uId := uuid.NewString()

	return models.User{
		ID:          uId,
		IdpSubject:  uId,
		IdpProvider: "google",
		Username:    "user" + uId,
		Email:       "user@test.com",
		IsActive:    true,
	}
}

func (f *Fixtures) InsertUser(u models.User) string {
	if f.db != nil {
		if err := f.db.WithContext(f.ctx).Create(&u).Error; err != nil {
			panic(fmt.Sprintf("insert user fixture failed: %v", err))
		}
		return u.ID
	}
	return ""
}
