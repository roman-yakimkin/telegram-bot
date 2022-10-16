package userstateprocessors

import "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"

const AmountRegex = `^(\d*)([\.,](\d{2})?)?$`

type UserStateProcessor interface {
	GetProcessStatus() int
	SetUserState(state *userstates.UserState)
	DoProcess(msgText string)
}
