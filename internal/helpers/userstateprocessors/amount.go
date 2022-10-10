package userstateprocessors

import (
	"regexp"
	"strconv"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

const AmountRegex = `^(\d*)([\.,](\d{2})?)?$`

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

func (p *amountProcessor) parseMsgText(msgText string) (int, int, error) {
	var err error
	matches := p.amountRegExp.FindAllStringSubmatch(msgText, -1)
	if len(matches) != 1 {
		return 0, 0, localerr.ErrIncorrectAmountValue
	}
	intPartStr, fracPartStr := matches[0][1], matches[0][3]
	if intPartStr == "" && fracPartStr == "" {
		return 0, 0, localerr.ErrIncorrectAmountValue
	}
	var intPart, fracPart int
	if intPartStr != "" {
		intPart, err = strconv.Atoi(intPartStr)
		if err != nil {
			return 0, 0, err
		}
	}
	if fracPartStr != "" {
		fracPart, err = strconv.Atoi(fracPartStr)
		if err != nil {
			return 0, 0, err
		}
	}
	return intPart, fracPart, nil
}

func (p *amountProcessor) DoProcess(msgText string) {
	amountInt, amountFrac, err := p.parseMsgText(msgText)
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectAmount)
		return
	}
	amount := amountInt*100 + amountFrac
	p.userState.SetUnconvertedAmount(amount)
}
