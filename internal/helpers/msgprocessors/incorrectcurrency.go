package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type incorrectCurrencyMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewIncorrectCurrencyMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &incorrectCurrencyMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *incorrectCurrencyMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.IncorrectCurrency
}

func (p *incorrectCurrencyMessageProcessor) DoProcess(ctx context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	messageId := "setcurrency_incorrectcurrency"
	currOutput, err := p.output.Currency().Output(ctx)
	if err != nil {
		return userstates.ExpectedCommand, messageId, err
	}
	return userstates.ExpectedCurrency, messageId, p.tgClient.SendMessage("Валюта задана неверно\n"+currOutput, msg.UserId)
}
