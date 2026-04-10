package repositories

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/ptracker/internal/apierr"
	"github.com/ptracker/internal/domain"
	"github.com/ptracker/internal/interfaces"
	"github.com/ptracker/internal/repositories/models"
	"gorm.io/gorm"
)

type UserRepo struct {
	db *gorm.DB
}

func NewUserRepo(db *gorm.DB) interfaces.UserRepository {
	return &UserRepo{
		db: db,
	}
}

func (r *UserRepo) Create(ctx context.Context,
	idpSubject, idpProvider, username, email string,
	displayName, avatarUrl *string) (string, error) {

	id := uuid.NewString()
	user := models.User{
		ID:          id,
		IdpSubject:  idpSubject,
		IdpProvider: idpProvider,
		Username:    username,
		Email:       email,
		DisplayName: displayName,
		AvatarURL:   avatarUrl,
	}
	err := gorm.G[models.User](r.db).Create(ctx, &user)
	if err != nil {
		return "", fmt.Errorf("store create user: %w", err)
	}

	return user.ID, nil
}

func (r *UserRepo) Get(ctx context.Context, id string) (domain.User, error) {

	user, err := gorm.G[models.User](r.db).Where("id = ?", id).First(ctx)
	if err == gorm.ErrRecordNotFound {
		return user.ToUserDomain(), apierr.ErrNotFound
	} else if err != nil {
		return user.ToUserDomain(), fmt.Errorf("gorm query: %w", err)
	}

	return user.ToUserDomain(), nil
}

func (r *UserRepo) GetBySubject(ctx context.Context,
	idpSubject, idpProvider string) (domain.User, error) {

	user, err := gorm.G[models.User](r.db).Where(
		"idp_subject = ? AND idp_provider = ?",
		idpSubject, idpProvider).First(ctx)
	if err == gorm.ErrRecordNotFound {
		return user.ToUserDomain(), apierr.ErrNotFound
	} else if err != nil {
		return user.ToUserDomain(), fmt.Errorf("gorm query: %w", err)
	}

	return user.ToUserDomain(), nil
}
