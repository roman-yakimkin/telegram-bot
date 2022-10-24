package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type startMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewStartMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &startMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *startMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/start"
}

func (p *startMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, error) {
	return userstates.ExpectedCommand, p.tgClient.SendMessage("hello\n"+InfoText, msg.UserId)
}
