package userstateprocessors

import (
	"strconv"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type delLimitMonthProcessor struct {
	processStatus int
	userState     *userstates.UserState
}

func NewDelLimitMonthProcessor() UserStateProcessor {
	return &delLimitMonthProcessor{
		processStatus: userstates.ExpectedDelLimitMonth,
	}
}

func (p *delLimitMonthProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *delLimitMonthProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *delLimitMonthProcessor) DoProcess(msgText string) {
	if msgText == "*" {
		p.userState.SetStatus(userstates.ExpectedCommand)
		return
	}
	index, err := strconv.Atoi(msgText)
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectDelLimitMonth)
		return
	}
	if index < 1 || index > 12 {
		p.userState.SetStatus(userstates.IncorrectDelLimitMonth)
		return
	}
	p.userState.SetBufferValue(userstates.DeleteLimitMonthIndex, index)
}
