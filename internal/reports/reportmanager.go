package reports

import (
	"fmt"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"strings"
	"time"
)

type ReportManager struct {
	er repo.ExpensesRepo
}

func New(er repo.ExpensesRepo) *ReportManager {
	return &ReportManager{
		er: er,
	}
}

func (rm *ReportManager) makeTextReport(expData map[string]int) string {
	var sb strings.Builder
	for cat, amount := range expData {
		sb.WriteString(fmt.Sprintf("Категория: %s, Трата: %d₽\n", cat, amount))
	}
	if len(expData) == 0 {
		sb.WriteString("Нет информации о расходах в данных период времени")
	}
	return sb.String()
}

func (rm *ReportManager) LastWeek(UserID int64) (result string) {
	timeStart := time.Now().AddDate(0, 0, -7)
	timeEnd := time.Now()
	expData := rm.er.ExpensesByUserAndTimeInterval(UserID, timeStart, timeEnd)
	return rm.makeTextReport(expData)
}

func (rm *ReportManager) LastMonth(UserID int64) string {
	timeStart := time.Now().AddDate(0, -1, 0)
	timeEnd := time.Now()
	expData := rm.er.ExpensesByUserAndTimeInterval(UserID, timeStart, timeEnd)
	return rm.makeTextReport(expData)
}

func (rm *ReportManager) LastYear(UserID int64) string {
	timeStart := time.Now().AddDate(-1, 0, 0)
	timeEnd := time.Now()
	expData := rm.er.ExpensesByUserAndTimeInterval(UserID, timeStart, timeEnd)
	return rm.makeTextReport(expData)
}
