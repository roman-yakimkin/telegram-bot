package msgprocessors

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type incorrectDateMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewIncorrectDateMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &incorrectDateMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *incorrectDateMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.IncorrectDate
}

func (p *incorrectDateMessageProcessor) DoProcess(msg Message, _ *userstates.UserState) (int, error) {
	return userstates.ExpectedDate, p.tgClient.SendMessage("Дата задана некорректно. Введите дату платежа в формате ГГГГ-ММ-ДД (* - текущая дата)", msg.UserID)
}
