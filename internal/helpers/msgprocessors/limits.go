package msgprocessors

import (
	"log"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type limitsMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewLimitsMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &limitsMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *limitsMessageProcessor) ShouldProcess(msg Message, _ *userstates.UserState) bool {
	return msg.Text == "/limits"
}

func (p *limitsMessageProcessor) DoProcess(msg Message, _ *userstates.UserState) (int, error) {
	limits, err := p.output.Limits().Output(msg.UserID)
	if err != nil {
		limits = "Ошибка при выводе лимитов"
		log.Println("limit output error", err)
	}
	return userstates.ExpectedCommand, p.tgClient.SendMessage(limits, msg.UserID)
}
