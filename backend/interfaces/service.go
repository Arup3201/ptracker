package interfaces

import (
	"context"

	"github.com/ptracker/domain"
)

type AuthService interface {
	RedirectLogin(ctx context.Context) (string, error)
	Callback(ctx context.Context,
		state, code string,
		userAgent, device, ipAddress string) (*domain.Session, error)
	Refresh(ctx context.Context,
		sessionId string) error
	Logout(ctx context.Context,
		sessionId string) error
}
