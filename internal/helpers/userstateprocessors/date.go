package userstateprocessors

import (
	"log"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store"
)

type DateProcessor struct {
	processStatus int
	userState     *userstates.UserState
	store         store.Store
	currConv      convertors.CurrencyConvertor
}

func NewDateProcessor(store store.Store, currConv convertors.CurrencyConvertor) UserStateProcessor {
	return &DateProcessor{
		processStatus: userstates.ExpectedDate,
		store:         store,
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
	} else {
		var err error
		date, err = time.Parse("2006-01-02", msgText)
		if err != nil {
			p.userState.SetStatus(userstates.IncorrectDate)
			return
		}
	}
	amountInBaseCurrency, err := p.convertAndAddAmount(date)
	if err != nil {
		log.Println("error on currency converting:", err)
		return
	}

	ok, err := p.checkLimitExceeding(amountInBaseCurrency, date)
	if err != nil {
		log.Println("error on limit exceeding checking:", err)
		return
	}
	if !ok {
		p.userState.SetStatus(userstates.LimitExceededAmount)
		return
	}

	p.userState.SetBufferValue(userstates.AddExpenseDateValue, date)
}

func (p *DateProcessor) convertAndAddAmount(date time.Time) (int, error) {
	amount, err := p.userState.IfFloatTransformToInt(userstates.AddExpenseAmountValue)
	if err != nil {
		log.Println("error upon getting expense amount", err)
		return 0, err
	}
	amountInBaseCurrency, err := p.currConv.From(amount, p.userState.Currency, utils.TimeTruncate(date))
	if err != nil {
		return 0, err
	}
	p.userState.SetBufferValue(userstates.AddExpenseAmountValue, amountInBaseCurrency)
	return amountInBaseCurrency, nil
}

func (p *DateProcessor) checkLimitExceeding(amount int, date time.Time) (bool, error) {
	return p.store.MeetMonthlyLimit(p.userState.UserID, utils.TimeTruncate(date), amount, p.currConv)
}
