package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type expectedDateMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewExpectedDateMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &expectedDateMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *expectedDateMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.ExpectedDate
}

func (p *expectedDateMessageProcessor) DoProcess(_ context.Context, msg Message, userState *userstates.UserState) (int, error) {
	return userstates.ExpectedCommand, p.tgClient.SendMessage("Информация о платеже добавлена", msg.UserId)
}
