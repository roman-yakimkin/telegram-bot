package msgprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/output"
)

type incorrectDelLimitMonthMessageProcessor struct {
	tgClient MessageSender
	output   output.Output
}

func NewIncorrectDelLimitMonthMessageProcessor(ms MessageSender, output output.Output) MessageProcessor {
	return &incorrectDelLimitMonthMessageProcessor{
		tgClient: ms,
		output:   output,
	}
}

func (p *incorrectDelLimitMonthMessageProcessor) ShouldProcess(_ Message, userState *userstates.UserState) bool {
	return userState.GetStatus() == userstates.IncorrectDelLimitMonth
}

func (p *incorrectDelLimitMonthMessageProcessor) DoProcess(_ context.Context, msg Message, _ *userstates.UserState) (int, string, error) {
	return userstates.ExpectedDelLimitMonth, "dellimit_incorrectmonth", p.tgClient.SendMessage("Месяц задан неверно. Введите месяц (1 - 12) или * для отмены", msg.UserId)
}
