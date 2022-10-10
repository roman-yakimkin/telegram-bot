package userstateprocessors

import (
	"log"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
)

type DateProcessor struct {
	processStatus int
	userState     *userstates.UserState
	currConv      convertors.CurrencyConvertorFrom
}

func NewDateProcessor(currConv convertors.CurrencyConvertorFrom) UserStateProcessor {
	return &DateProcessor{
		processStatus: userstates.ExpectedDate,
		currConv:      currConv,
	}
}

func (p *DateProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *DateProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *DateProcessor) DoProcess(msgText string) {
	var err error
	var date time.Time
	if msgText == "*" {
		date = time.Now()
		p.userState.SetDate(date)
	} else {
		var err error
		date, err = time.Parse("2006-01-02", msgText)
		if err != nil {
			p.userState.SetStatus(userstates.IncorrectDate)
			return
		}
	}
	if err = p.convertAndAddAmount(date); err != nil {
		log.Println("error on currency converting:", err)
		return
	}
	p.userState.SetDate(date)
}

func (p *DateProcessor) convertAndAddAmount(date time.Time) error {
	amount := p.userState.GetAmount()
	amountInBaseCurrency, err := p.currConv.From(amount, p.userState.Currency, utils.TimeTruncate(date))
	if err != nil {
		return err
	}
	p.userState.SetConvertedAmount(amountInBaseCurrency)
	return nil
}
