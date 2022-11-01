package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type incorrectCategoryMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewIncorrectCategoryMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &incorrectCategoryMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *incorrectCategoryMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.IncorrectCategory
}

func (p *incorrectCategoryMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	return userstates.ExpectedCategory, MessageNewExpenseIncorrectCategory, p.tgClient.SendMessage("Категория задана неверно. Введите категорию платежа", msg.UserId)
}
