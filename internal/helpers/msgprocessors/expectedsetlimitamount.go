package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type expectedSetLimitAmountMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewExpectedSetLimitAmountMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &expectedSetLimitAmountMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *expectedSetLimitAmountMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.ExpectedSetLimitAmount
}

func (p *expectedSetLimitAmountMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, error) {
	return userstates.ExpectedCommand, p.tgClient.SendMessage("Лимит установлен", msg.UserID)
}
