package userstateprocessors

import (
	"strconv"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type setLimitMonthProcessor struct {
	processStatus int
	userState     *userstates.UserState
}

func NewSetLimitMonthProcessor() UserStateProcessor {
	return &setLimitMonthProcessor{
		processStatus: userstates.ExpectedSetLimitMonth,
	}
}

func (p *setLimitMonthProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *setLimitMonthProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *setLimitMonthProcessor) DoProcess(msgText string) {
	if msgText == "*" {
		p.userState.SetStatus(userstates.ExpectedCommand)
		return
	}
	limit, err := strconv.Atoi(msgText)
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectSetLimitMonth)
		return
	}
	if limit < 1 || limit > 12 {
		p.userState.SetStatus(userstates.IncorrectSetLimitMonth)
		return
	}
	p.userState.SetBufferValue(userstates.SetLimitMonthIndex, limit)
}
