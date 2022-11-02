package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type incorrectDateMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewIncorrectDateMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &incorrectDateMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *incorrectDateMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.IncorrectDate
}

func (p *incorrectDateMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	return userstates.ExpectedDate, MessageNewExpenseIncorrectDate, p.tgClient.SendMessage("Дата задана некорректно. Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserId)
}
