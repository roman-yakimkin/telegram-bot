package reports

import (
	"fmt"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/interfaces"
	"time"
)

type ReportManager struct {
	er interfaces.ExpensesRepo
}

func New(er interfaces.ExpensesRepo) *ReportManager {
	return &ReportManager{
		er: er,
	}
}

func (rm *ReportManager) makeTextReport(expData map[string]int) string {
	var result string
	for cat, amount := range expData {
		result += fmt.Sprintf("%s : %d\n", cat, amount)
	}
	if result == "" {
		result = "Нет информации о расходах в данных период времени"
	}
	return result
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
