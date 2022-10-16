package msgprocessors

import (
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

func (p *expectedDelLimitMonthMessageProcessor) DoProcess(msg Message, userState *userstates.UserState) (int, error) {
	return userstates.ExpectedCommand, p.tgClient.SendMessage("Лимит удален до состояния по умолчанию", msg.UserID)
}
