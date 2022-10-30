package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type expectedSetLimitMonthMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewExpectedSetLimitMonthMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &expectedSetLimitMonthMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *expectedSetLimitMonthMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.ExpectedSetLimitMonth
}

func (p *expectedSetLimitMonthMessageProcessor) DoProcess(_ context.Context, msg Message, userState *userstates.UserState) (int, string, error) {
	return userstates.ExpectedSetLimitAmount, "setlimit_month", p.tgClient.SendMessage("Введите сумму лимита. Текущая валюта - "+userState.Currency, msg.UserId)
}
