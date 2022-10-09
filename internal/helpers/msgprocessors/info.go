package msgprocessors

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type infoMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewInfoMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &infoMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *infoMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/info"
}

func (p *infoMessageProcessor) DoProcess(msg Message, _ *userstates.UserState) (int, error) {
	return userstates.ExpectedCommand, p.tgClient.SendMessage(InfoText, msg.UserID)
}
