package repo

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type UserStateRepo interface {
	GetOne(ctx context.Context, UserId int64) (*userstates.UserState, error)
	Save(ctx context.Context, state *userstates.UserState) error
	Delete(ctx context.Context, UserID int64) error
	ClearStatus(ctx context.Context) error
}
