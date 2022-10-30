package msgprocessors

import (
	"context"

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

func (p *incorrectSetLimitMonthMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	return userstates.ExpectedSetLimitMonth, "setlimit_incorrectmonth", p.tgClient.SendMessage("Месяц задан неверно. Введите месяц (1 - 12)", msg.UserId)
}
