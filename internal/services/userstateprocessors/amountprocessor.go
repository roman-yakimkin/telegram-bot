package userstateprocessors

import (
	"strconv"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type AmountProcessor struct {
	processStatus int
	userState     *userstates.UserState
}

func NewAmountProcessor() *AmountProcessor {
	return &AmountProcessor{
		processStatus: userstates.ExpectedAmount,
	}
}

func (p *AmountProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *AmountProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *AmountProcessor) DoProcess(msgText string) {
	var err error
	amount, err := strconv.Atoi(msgText)
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectAmount)
		return
	}
	if amount < 0 {
		p.userState.SetStatus(userstates.IncorrectAmount)
		return
	}
	p.userState.SetAmount(amount)
}
