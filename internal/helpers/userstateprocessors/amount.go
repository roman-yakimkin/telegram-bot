package userstateprocessors

import (
	"context"
	"regexp"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type amountProcessor struct {
	processStatus int
	amountRegExp  *regexp.Regexp
}

func NewAmountProcessor() (UserStateProcessor, error) {
	amountRegExp, err := regexp.Compile(AmountRegex)
	if err != nil {
		return nil, err
	}
	return &amountProcessor{
		processStatus: userstates.ExpectedAmount,
		amountRegExp:  amountRegExp,
	}, nil
}

func (p *amountProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *amountProcessor) DoProcess(_ context.Context, state *userstates.UserState, msgText string) {
	amountInt, amountFrac, err := utils.ParseMsgText(msgText, p.amountRegExp)
	if err != nil {
		state.SetStatus(userstates.IncorrectAmount)
		return
	}
	amount := amountInt*100 + amountFrac
	state.SetBufferValue(userstates.AddExpenseAmountValue, amount)
}
