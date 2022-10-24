package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type expectedCurrencyMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewExpectedCurrencyMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &expectedCurrencyMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *expectedCurrencyMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.ExpectedCurrency
}

func (p *expectedCurrencyMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, error) {
	return userstates.ExpectedCommand, p.tgClient.SendMessage("Валюта изменена", msg.UserId)
}
