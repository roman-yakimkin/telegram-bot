package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type deleteLimitMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewDeleteLimitMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &deleteLimitMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *deleteLimitMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/dellimit"
}

func (p *deleteLimitMessageProcessor) DoProcess(_ context.Context, msg Message, userState *userstates.UserState) (int, string, error) {
	return userstates.ExpectedDelLimitMonth, MessageDelLimit, p.tgClient.SendMessage("Введите месяц (1 - 12) или * для отмены", msg.UserId)
}
