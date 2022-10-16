package msgprocessors

import (
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type incorrectSetLimitMonthMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewIncorrectSetLimitMonthMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &incorrectSetLimitMonthMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *incorrectSetLimitMonthMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.IncorrectSetLimitMonth
}

func (p *incorrectSetLimitMonthMessageProcessor) DoProcess(msg Message, _ *userstates.UserState) (int, error) {
	return userstates.ExpectedSetLimitMonth, p.tgClient.SendMessage("Месяц задан неверно. Введите месяц (1 - 12)", msg.UserID)
}
