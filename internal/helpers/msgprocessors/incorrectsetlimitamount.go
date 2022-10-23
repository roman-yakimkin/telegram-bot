package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type incorrectSetLimitAmountMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewIncorrectSetLimitAmountMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &incorrectSetLimitAmountMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *incorrectSetLimitAmountMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.IncorrectSetLimitAmount
}

func (p *incorrectSetLimitAmountMessageProcessor) DoProcess(_ context.Context, msg Message, userState *userstates.UserState) (int, error) {
	return userstates.ExpectedSetLimitAmount, p.tgClient.SendMessage("Сумма лимита задана неверно. Введите сумму лимита. Текущая валюта - "+userState.Currency, msg.UserID)
}
