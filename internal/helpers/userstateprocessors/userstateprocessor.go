package userstateprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

const AmountRegex = `^(\d*)([\.,](\d{2})?)?$`

type UserStateProcessor interface {
	GetProcessStatus() int
	DoProcess(ctx context.Context, state *userstates.UserState, msgText string)
}
