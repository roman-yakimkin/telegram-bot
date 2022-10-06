package userstateprocessors

import (
	"log"
	"regexp"
	"strconv"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/localerr"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

const AmountRegex = `^(\d*)([\.,](\d{2})?)?$`

type AmountProcessor struct {
	processStatus int
	userState     *userstates.UserState
	amountRegExp  *regexp.Regexp
	currConv      convertors.CurrencyConvertorFrom
}

func NewAmountProcessor(currConv convertors.CurrencyConvertorFrom) *AmountProcessor {
	return &AmountProcessor{
		processStatus: userstates.ExpectedAmount,
		currConv:      currConv,
	}
}

func (p *AmountProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *AmountProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *AmountProcessor) parseMsgText(msgText string) (int, int, error) {
	var err error
	if p.amountRegExp == nil {
		p.amountRegExp, err = regexp.Compile(AmountRegex)
		if err != nil {
			return 0, 0, err
		}
	}
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

func (p *AmountProcessor) DoProcess(msgText string) {
	amountInt, amountFrac, err := p.parseMsgText(msgText)
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectAmount)
		return
	}
	amount := amountInt*100 + amountFrac
	amountInBaseCurrency, err := p.currConv.From(amount, p.userState.Currency)
	if err != nil {
		log.Println("error on currency converting:", err)
		return
	}
	p.userState.SetAmount(amountInBaseCurrency)
}
