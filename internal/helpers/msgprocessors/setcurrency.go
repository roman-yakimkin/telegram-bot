package msgprocessors

import (
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

func (p *setCurrencyMessageProcessor) DoProcess(msg Message, _ *userstates.UserState) (int, error) {
	currOutput, err := p.output.Currency().Output()
	if err != nil {
		return userstates.ExpectedCommand, err
	}
	return userstates.ExpectedCurrency, p.tgClient.SendMessage(currOutput, msg.UserID)
}
