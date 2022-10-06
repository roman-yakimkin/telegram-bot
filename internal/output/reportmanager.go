package output

import (
	"fmt"
	"strings"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store"
)

type ReportManager struct {
	store      store.Store
	conv       convertors.CurrencyConvertorTo
	currAmount *CurrencyAmount
}

func NewReportManager(store store.Store, conv convertors.CurrencyConvertorTo, currAmount *CurrencyAmount) *ReportManager {
	return &ReportManager{
		store:      store,
		conv:       conv,
		currAmount: currAmount,
	}
}

func (rm *ReportManager) makeTextReport(userID int64, expData repo.ExpData) (string, error) {
	var sb strings.Builder
	for cat, expByCurrency := range expData {
		amountStrs := make([]string, 0, len(expByCurrency))
		for currName, amount := range expByCurrency {
			amountInCurrency, err := rm.conv.To(amount, currName)
			if err != nil {
				return "", err
			}
			amountDisplay, err := rm.currAmount.Output(amountInCurrency, currName)
			if err != nil {
				return "", err
			}
			amountStrs = append(amountStrs, amountDisplay)
		}
		sb.WriteString(fmt.Sprintf("Категория: %s, Трата: %s\n", cat, strings.Join(amountStrs, ", ")))
	}
	if len(expData) == 0 {
		sb.WriteString("Нет информации о расходах в данных период времени")
	}
	return sb.String(), nil
}

func (rm *ReportManager) LastWeek(UserID int64) (string, error) {
	timeStart := time.Now().AddDate(0, 0, -7)
	timeEnd := time.Now()
	expData := rm.store.Expense().ExpensesByUserAndTimeInterval(UserID, timeStart, timeEnd)
	result, err := rm.makeTextReport(UserID, expData)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (rm *ReportManager) LastMonth(UserID int64) (string, error) {
	timeStart := time.Now().AddDate(0, -1, 0)
	timeEnd := time.Now()
	expData := rm.store.Expense().ExpensesByUserAndTimeInterval(UserID, timeStart, timeEnd)
	return rm.makeTextReport(UserID, expData)
}

func (rm *ReportManager) LastYear(UserID int64) (string, error) {
	timeStart := time.Now().AddDate(-1, 0, 0)
	timeEnd := time.Now()
	expData := rm.store.Expense().ExpensesByUserAndTimeInterval(UserID, timeStart, timeEnd)
	return rm.makeTextReport(UserID, expData)
}
