package userstateprocessors

import "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"

type UserStateProcessor interface {
	GetProcessStatus() int
	SetUserState(state *userstates.UserState)
	DoProcess(msgText string)
}
