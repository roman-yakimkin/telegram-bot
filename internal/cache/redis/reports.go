package redis

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/cache"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/defs"
)

type cacheReports struct {
	client *redis.Client
}

func NewCacheReports(client *redis.Client) cache.Reports {
	return &cacheReports{
		client: client,
	}
}

func (r *cacheReports) getCacheKey(userId int64, reportType defs.ReportType) string {
	return fmt.Sprintf("user:%d:type:%s", userId, reportType)
}

func (r *cacheReports) Get(userId int64, reportType defs.ReportType) (string, error) {
	key := r.getCacheKey(userId, reportType)
	return r.client.Get(key).Result()
}

func (r *cacheReports) Set(userId int64, reportType defs.ReportType, reportStr string, invalidateDate time.Time) error {
	key := r.getCacheKey(userId, reportType)
	duration := time.Until(invalidateDate)
	return r.client.Set(key, reportStr, duration).Err()
}

func (r *cacheReports) invalidateOneReport(userId int64, reportType defs.ReportType) error {
	key := r.getCacheKey(userId, reportType)
	return r.client.Del(key).Err()
}

func (r *cacheReports) Invalidate(userId int64, newExpenseDate time.Time) error {
	if newExpenseDate.After(time.Now().AddDate(-1, 0, 0)) {
		err := r.invalidateOneReport(userId, defs.ReportLastYear)
		if err != nil {
			return err
		}
	}
	if newExpenseDate.After(time.Now().AddDate(0, -1, 0)) {
		err := r.invalidateOneReport(userId, defs.ReportLastMonth)
		if err != nil {
			return err
		}
	}
	if newExpenseDate.After(time.Now().AddDate(0, 0, -7)) {
		err := r.invalidateOneReport(userId, defs.ReportLastWeek)
		if err != nil {
			return err
		}
	}
	return nil
}
