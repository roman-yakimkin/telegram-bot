package userstateprocessors

import (
	"context"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type CategoryProcessor struct {
	processStatus int
}

func NewCategoryProcessor() UserStateProcessor {
	return &CategoryProcessor{
		processStatus: userstates.ExpectedCategory,
	}
}

func (p *CategoryProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *CategoryProcessor) DoProcess(_ context.Context, state *userstates.UserState, msgText string) {
	state.SetBufferValue(userstates.AddExpenseCategoryValue, msgText)
}
