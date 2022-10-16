package userstateprocessors

import "gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"

type CategoryProcessor struct {
	processStatus int
	userState     *userstates.UserState
}

func NewCategoryProcessor() UserStateProcessor {
	return &CategoryProcessor{
		processStatus: userstates.ExpectedCategory,
	}
}

func (p *CategoryProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *CategoryProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *CategoryProcessor) DoProcess(msgText string) {
	p.userState.SetBufferValue(userstates.AddExpenseCategoryValue, msgText)
}
