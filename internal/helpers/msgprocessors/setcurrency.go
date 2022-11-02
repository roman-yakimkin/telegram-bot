package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type setCurrencyMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewSetCurrencyMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &setCurrencyMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *setCurrencyMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/setcurrency"
}

func (p *setCurrencyMessageProcessor) DoProcess(ctx context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	currOutput, err := p.output.Currency().Output(ctx)
	if err != nil {
		return userstates.ExpectedCommand, MessageSetCurrency, err
	}
	return userstates.ExpectedCurrency, MessageSetCurrency, p.tgClient.SendMessage(currOutput, msg.UserId)
}
