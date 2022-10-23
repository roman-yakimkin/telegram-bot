package userstateprocessors

import (
	"context"
	"strconv"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type delLimitMonthProcessor struct {
	processStatus int
}

func NewDelLimitMonthProcessor() UserStateProcessor {
	return &delLimitMonthProcessor{
		processStatus: userstates.ExpectedDelLimitMonth,
	}
}

func (p *delLimitMonthProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *delLimitMonthProcessor) DoProcess(_ context.Context, state *userstates.UserState, msgText string) {
	if msgText == "*" {
		state.SetStatus(userstates.ExpectedCommand)
		return
	}
	index, err := strconv.Atoi(msgText)
	if err != nil {
		state.SetStatus(userstates.IncorrectDelLimitMonth)
		return
	}
	if index < 1 || index > 12 {
		state.SetStatus(userstates.IncorrectDelLimitMonth)
		return
	}
	state.SetBufferValue(userstates.DeleteLimitMonthIndex, index)
}
