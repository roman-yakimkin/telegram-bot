package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type incorrectAmountMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewIncorrectAmountMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &incorrectAmountMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *incorrectAmountMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.IncorrectAmount
}

func (p *incorrectAmountMessageProcessor) DoProcess(_ context.Context, msg Message, userState *userstates.UserState) (int, error) {
	return userstates.ExpectedAmount, p.tgClient.SendMessage("Сумма платежа задана неверно. Введите сумму платежа. Текущая валюта - "+userState.Currency, msg.UserId)
}
