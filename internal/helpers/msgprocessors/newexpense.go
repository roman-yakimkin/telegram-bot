package msgprocessors

import (
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

func (p *newExpenseMessageProcessor) DoProcess(msg Message, _ *userstates.UserState) (int, error) {
	return userstates.ExpectedCategory, p.tgClient.SendMessage("Введите категорию платежа", msg.UserID)
}
