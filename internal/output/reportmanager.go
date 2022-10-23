package output

import (
	"context"
	"fmt"
	"log"
	"strings"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/convertors"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store"
)

type ReportManagerLastWeek interface {
	LastWeek(ctx context.Context, UserID int64) (string, error)
}

type ReportManagerLastMonth interface {
	LastMonth(ctx context.Context, UserID int64) (string, error)
}

type ReportManagerLastYear interface {
	LastYear(ctx context.Context, UserID int64) (string, error)
}

type ReportManager interface {
	ReportManagerLastWeek
	ReportManagerLastMonth
	ReportManagerLastYear
}

type reportManager struct {
	store      store.Store
	conv       convertors.CurrencyConvertorTo
	currAmount CurrencyAmount
}

func NewReportManager(store store.Store, conv convertors.CurrencyConvertorTo, currAmount CurrencyAmount) ReportManager {
	return &reportManager{
		store:      store,
		conv:       conv,
		currAmount: currAmount,
	}
}

func (rm *reportManager) makeTextReport(ctx context.Context, userID int64, expData repo.ExpData) (string, error) {
	var sb strings.Builder
	for cat, expByCurrency := range expData {
		amountStrs := make([]string, 0, len(expByCurrency))
		for currName, amountMap := range expByCurrency {
			var amountTotal int
			for date, amount := range amountMap {
				amountInCurrency, err := rm.conv.To(ctx, amount, currName, date)
				if err != nil {
					return "", err
				}
				amountTotal += amountInCurrency
			}
			amountDisplay, err := rm.currAmount.Output(ctx, amountTotal, currName)
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

func (rm *reportManager) LastWeek(ctx context.Context, UserID int64) (string, error) {
	timeStart := time.Now().AddDate(0, 0, -7)
	timeEnd := time.Now()
	expData, err := rm.store.Expense().ExpensesByUserAndTimeInterval(ctx, UserID, timeStart, timeEnd)
	if err != nil {
		log.Println("getting report data error:", err)
		return "Ошибка при получении данных", err
	}
	result, err := rm.makeTextReport(ctx, UserID, expData)
	if err != nil {
		return "", err
	}
	return result, nil
}

func (rm *reportManager) LastMonth(ctx context.Context, UserID int64) (string, error) {
	timeStart := time.Now().AddDate(0, -1, 0)
	timeEnd := time.Now()
	expData, err := rm.store.Expense().ExpensesByUserAndTimeInterval(ctx, UserID, timeStart, timeEnd)
	if err != nil {
		log.Println("getting report data error:", err)
		return "Ошибка при получении данных", err
	}
	return rm.makeTextReport(ctx, UserID, expData)
}

func (rm *reportManager) LastYear(ctx context.Context, UserID int64) (string, error) {
	timeStart := time.Now().AddDate(-1, 0, 0)
	timeEnd := time.Now()
	expData, err := rm.store.Expense().ExpensesByUserAndTimeInterval(ctx, UserID, timeStart, timeEnd)
	if err != nil {
		log.Println("getting report data error:", err)
		return "Ошибка при получении данных", err
	}
	return rm.makeTextReport(ctx, UserID, expData)
}
