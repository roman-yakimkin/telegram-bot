package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type newExpenseMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewNewExpenseMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &newExpenseMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *newExpenseMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/newexpense"
}

func (p *newExpenseMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	return userstates.ExpectedCategory, "newexpense", p.tgClient.SendMessage("Введите категорию платежа", msg.UserId)
}
