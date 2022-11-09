package cache

import (
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/defs"
)

type Reports interface {
	Get(userId int64, reportType defs.ReportType) (string, error)
	Set(userId int64, reportType defs.ReportType, reportStr string, invalidateDate time.Time) error
	Invalidate(userId int64, newExpenseDate time.Time) error
}
