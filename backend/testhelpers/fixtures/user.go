package fixtures

import (
	"fmt"

	"github.com/google/uuid"
	"github.com/ptracker/models"
)

func RandomUserRow() models.User {
	uId := uuid.NewString()

	return models.User{
		ID:       uId,
		Username: "user" + uId,
		Email:    "user@test.com",
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
