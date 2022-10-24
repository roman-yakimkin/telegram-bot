package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type setLimitMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewSetLimitMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &setLimitMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *setLimitMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/setlimit"
}

func (p *setLimitMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, error) {
	return userstates.ExpectedSetLimitMonth, p.tgClient.SendMessage("Введите месяц (1 - 12) или * для отмены", msg.UserId)
}
