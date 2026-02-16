package interfaces

import (
	"context"

	"github.com/ptracker/internal/domain"
)

type Notifier interface {
	Notify(ctx context.Context, user string, message domain.Message) error
	BatchNotify(ctx context.Context, users []string, message domain.Message) error
}
