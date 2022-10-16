package userstateprocessors

import (
	"regexp"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type amountProcessor struct {
	processStatus int
	userState     *userstates.UserState
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

func (p *amountProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *amountProcessor) DoProcess(msgText string) {
	amountInt, amountFrac, err := utils.ParseMsgText(msgText, p.amountRegExp)
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectAmount)
		return
	}
	amount := amountInt*100 + amountFrac
	p.userState.SetBufferValue(userstates.AddExpenseAmountValue, amount)
}
