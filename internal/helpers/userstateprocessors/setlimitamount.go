package userstateprocessors

import (
	"log"
	"regexp"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type setLimitAmountProcessor struct {
	processStatus int
	userState     *userstates.UserState
	amountRegExp  *regexp.Regexp
	currConv      convertors.CurrencyConvertorFrom
}

func NewSetLimitAmountProcessor(currConv convertors.CurrencyConvertorFrom) (UserStateProcessor, error) {
	amountRegExp, err := regexp.Compile(AmountRegex)
	if err != nil {
		return nil, err
	}
	return &setLimitAmountProcessor{
		processStatus: userstates.ExpectedSetLimitAmount,
		amountRegExp:  amountRegExp,
		currConv:      currConv,
	}, nil
}

func (p *setLimitAmountProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *setLimitAmountProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *setLimitAmountProcessor) DoProcess(msgText string) {
	amountInt, amountFrac, err := utils.ParseMsgText(msgText, p.amountRegExp)
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectSetLimitAmount)
		return
	}
	amount := amountInt*100 + amountFrac
	if err = p.convertAndAddAmount(amount); err != nil {
		log.Println("error on currency converting:", err)
		p.userState.SetStatus(userstates.ExpectedCommand)
		return
	}
}

func (p *setLimitAmountProcessor) convertAndAddAmount(unconvertedAmount int) error {
	amountInBaseCurrency, err := p.currConv.From(unconvertedAmount, p.userState.Currency, utils.TimeTruncate(time.Now()))
	if err != nil {
		return err
	}
	p.userState.SetBufferValue(userstates.SetLimitMonthValue, amountInBaseCurrency)
	return nil
}
