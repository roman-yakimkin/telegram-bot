package userstateprocessors

import (
	"context"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/model/userstates"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store"
	"go.uber.org/zap"
)

type DateProcessor struct {
	processStatus int
	store         store.Store
	currConv      convertors.CurrencyConvertor
	logger        *zap.Logger
}

func NewDateProcessor(store store.Store, currConv convertors.CurrencyConvertor, logger *zap.Logger) UserStateProcessor {
	return &DateProcessor{
		processStatus: userstates.ExpectedDate,
		store:         store,
		currConv:      currConv,
		logger:        logger,
	}
}

func (p *DateProcessor) GetProcessStatus() int {
	return p.processStatus
}

func (p *DateProcessor) DoProcess(ctx context.Context, state *userstates.UserState, msgText string) {
	var err error
	var date time.Time
	if msgText == "*" {
		date = time.Now()
	} else {
		var err error
		date, err = time.Parse("2006-01-02", msgText)
		if err != nil {
			state.SetStatus(userstates.IncorrectDate)
			return
		}
	}
	amountInBaseCurrency, err := p.convertAndAddAmount(ctx, state, date)
	if err != nil {
		p.logger.Error("error on currency converting:", zap.Error(err))
		return
	}

	ok, err := p.checkLimitExceeding(ctx, state, amountInBaseCurrency, date)
	if err != nil {
		p.logger.Error("error on limit exceeding checking:", zap.Error(err))
		return
	}
	if !ok {
		state.SetStatus(userstates.LimitExceededAmount)
		return
	}

	state.SetBufferValue(userstates.AddExpenseDateValue, date)
}

func (p *DateProcessor) convertAndAddAmount(ctx context.Context, state *userstates.UserState, date time.Time) (int, error) {
	amount, err := state.IfFloatTransformToInt(userstates.AddExpenseAmountValue)
	if err != nil {
		p.logger.Error("error upon getting expense amount", zap.Error(err))
		return 0, err
	}
	amountInBaseCurrency, err := p.currConv.From(ctx, amount, state.Currency, utils.TimeTruncate(date))
	if err != nil {
		return 0, err
	}
	state.SetBufferValue(userstates.AddExpenseAmountValue, amountInBaseCurrency)
	return amountInBaseCurrency, nil
}

func (p *DateProcessor) checkLimitExceeding(ctx context.Context, state *userstates.UserState, amount int, date time.Time) (bool, error) {
	return p.store.MeetMonthlyLimit(ctx, state.UserId, utils.TimeTruncate(date), amount, p.currConv)
}
