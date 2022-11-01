package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type getCurrencyMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewGetCurrencyMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &getCurrencyMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *getCurrencyMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/getcurrency"
}

func (p *getCurrencyMessageProcessor) DoProcess(_ context.Context, msg Message, userState *userstates.UserState) (int, string, error) {
	return userstates.ExpectedCommand, MessageGetCurrency, p.tgClient.SendMessage("Ваша текущая валюта - "+userState.Currency, msg.UserId)
}
