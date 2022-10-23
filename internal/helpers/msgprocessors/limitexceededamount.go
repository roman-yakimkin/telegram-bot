package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type limitExceededAmountMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewLimitExceededAmountMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &limitExceededAmountMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *limitExceededAmountMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.LimitExceededAmount
}

func (p *limitExceededAmountMessageProcessor) DoProcess(_ context.Context, msg Message, userState *userstates.UserState) (int, error) {
	return userstates.ExpectedAmount, p.tgClient.SendMessage("При данной сумме платежа возникнет превышение месячного лимита. Введите другую сумму или дату. Текущая валюта - "+userState.Currency, msg.UserID)
}
