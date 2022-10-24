package repoupdaters

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type UserStateRepoUpdater interface {
	ReadyToUpdate(state *userstates.UserState) bool
	UpdateRepo(ctx context.Context, state *userstates.UserState) error
	ClearData(state *userstates.UserState)
}
