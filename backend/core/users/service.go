package users

import (
	"context"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/ptracker/core"
)

type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	DisplayName   *string   `json:"display_name"` // nullable
	Email         string    `json:"email"`
	AvatarURL     *string   `json:"avatar_url"` // nullable
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
	LastLoginTime time.Time `json:"-"`
}

var emailRegex, _ = regexp.Compile(
	`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`,
)

type UserService struct {
	userRepo *UserRepository
}

func NewUserService(userRepo *UserRepository) *UserService {
	return &UserService{
		userRepo: userRepo,
	}
}

func (s *UserService) Create(ctx context.Context,
	username, email string,
	displayName, avatarUrl *string) (string, error) {

	if strings.Trim(email, " ") == "" {
		return "", fmt.Errorf("email is empty: %w", core.ErrInvalidValue)
	}

	if match := emailRegex.Find([]byte(email)); match == nil {
		return "", fmt.Errorf("email pattern not supported: %w", core.ErrInvalidValue)
	}

	if strings.Trim(username, " ") == "" {
		return "", fmt.Errorf("username is empty: %w", core.ErrInvalidValue)
	}

	id, err := s.userRepo.Create(ctx, username, email, displayName, avatarUrl)
	if err != nil {
		return "", fmt.Errorf("user repository create: %w", err)
	}

	return id, nil
}

func (s *UserService) Get(ctx context.Context,
	id string) (*User, error) {

	user, err := s.userRepo.Get(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("user repository get: %w", err)
	}

	return &User{
		ID:            user.ID,
		Username:      user.Username,
		Email:         user.Email,
		DisplayName:   user.DisplayName,
		AvatarURL:     user.AvatarURL,
		CreatedAt:     user.CreatedAt,
		UpdatedAt:     user.UpdatedAt,
		LastLoginTime: user.LastLoginTime,
	}, nil
}
