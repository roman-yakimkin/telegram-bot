package userstateprocessors

import (
	"context"
	"strconv"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type setLimitMonthProcessor struct {
	processStatus int
}

func NewSetLimitMonthProcessor() UserStateProcessor {
	return &setLimitMonthProcessor{
		processStatus: userstates.ExpectedSetLimitMonth,
	}
}

func (p *setLimitMonthProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *setLimitMonthProcessor) DoProcess(_ context.Context, state *userstates.UserState, msgText string) {
	if msgText == "*" {
		state.SetStatus(userstates.ExpectedCommand)
		return
	}
	limit, err := strconv.Atoi(msgText)
	if err != nil {
		state.SetStatus(userstates.IncorrectSetLimitMonth)
		return
	}
	if limit < 1 || limit > 12 {
		state.SetStatus(userstates.IncorrectSetLimitMonth)
		return
	}
	state.SetBufferValue(userstates.SetLimitMonthIndex, limit)
}
