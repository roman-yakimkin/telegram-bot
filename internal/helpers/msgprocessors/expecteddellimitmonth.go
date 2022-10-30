package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type expectedDelLimitMonthMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewExpectedDelLimitMonthMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &expectedDelLimitMonthMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *expectedDelLimitMonthMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.ExpectedDelLimitMonth
}

func (p *expectedDelLimitMonthMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	return userstates.ExpectedCommand, "dellimit_month", p.tgClient.SendMessage("Лимит удален до состояния по умолчанию", msg.UserId)
}
