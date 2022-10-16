package repoupdaters

import "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"

type UserStateRepoUpdater interface {
	SetUserState(userState *userstates.UserState)
	ReadyToUpdate() bool
	UpdateRepo() error
	ClearData()
}
