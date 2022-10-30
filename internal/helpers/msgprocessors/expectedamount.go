package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type expectedAmountMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewExpectedAmountMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &expectedAmountMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *expectedAmountMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.ExpectedAmount
}

func (p *expectedAmountMessageProcessor) DoProcess(_ context.Context, msg Message, userState *userstates.UserState) (int, string, error) {
	return userstates.ExpectedDate, "newexpense_amount", p.tgClient.SendMessage("Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserId)
}
