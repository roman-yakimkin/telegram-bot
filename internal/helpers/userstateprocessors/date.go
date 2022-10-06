package userstateprocessors

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type DateProcessor struct {
	processStatus int
	userState     *userstates.UserState
}

func NewDateProcessor() *DateProcessor {
	return &DateProcessor{
		processStatus: userstates.ExpectedDate,
	}
}

func (p *DateProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *DateProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *DateProcessor) DoProcess(msgText string) {
	if msgText == "*" {
		p.userState.SetDate(time.Now())
		return
	}
	var err error
	date, err := time.Parse("2006-01-02", msgText)
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectDate)
		return
	}
	p.userState.SetDate(date)
}
