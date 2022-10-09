package msgprocessors

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type expectedCategoryMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewExpectedCategoryMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &expectedCategoryMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *expectedCategoryMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.ExpectedCategory
}

func (p *expectedCategoryMessageProcessor) DoProcess(msg Message, userState *userstates.UserState) (int, error) {
	return userstates.ExpectedAmount, p.tgClient.SendMessage("Введите сумму платежа. Текущая валюта - "+userState.Currency, msg.UserID)
}
