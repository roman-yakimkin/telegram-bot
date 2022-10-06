package userstateprocessors

import (
	"strings"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type CurrencyProcessor struct {
	processStatus int
	userState     *userstates.UserState
	currRepo      repo.CurrencyRepo
}

func NewCurrencyProcessor(currRepo repo.CurrencyRepo) *CurrencyProcessor {
	return &CurrencyProcessor{
		processStatus: userstates.ExpectedCurrency,
		currRepo:      currRepo,
	}
}

func (p *CurrencyProcessor) SetUserState(userState *userstates.UserState) {
	p.userState = userState
}

func (p *CurrencyProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *CurrencyProcessor) DoProcess(msgText string) {
	if msgText == "*" {
		p.userState.SetStatus(userstates.ExpectedCommand)
		return
	}
	currencies, err := p.currRepo.GetAll()
	if err != nil {
		p.userState.SetStatus(userstates.IncorrectCurrency)
		return
	}
	for _, currency := range currencies {
		if strings.EqualFold(currency.Name, msgText) {
			p.userState.Currency = strings.ToUpper(currency.Name)
			return
		}
	}
	p.userState.SetStatus(userstates.IncorrectCurrency)
}
