package output

import (
	"context"
	"errors"
	"time"

	"github.com/go-redis/redis"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/cache"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/defs"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/store"
	"go.uber.org/zap"
)

type cachedReportManager struct {
	rm     ReportManager
	store  store.Store
	cache  cache.Reports
	logger *zap.Logger
}

func NewCachedReportManager(rm ReportManager, store store.Store, cache cache.Reports, logger *zap.Logger) ReportManager {
	return &cachedReportManager{
		rm:     rm,
		store:  store,
		cache:  cache,
		logger: logger,
	}
}

func (c *cachedReportManager) earliestDateSince(ctx context.Context, date time.Time) (time.Time, error) {
	return c.store.Expense().EarliestDateSince(ctx, date)
}

func (c *cachedReportManager) StartTimeByReport(reportType defs.ReportType) time.Time {
	return c.rm.StartTimeByReport(reportType)
}

type reportGetFunc func(context.Context, int64) (string, error)

func (c *cachedReportManager) lastPeriod(ctx context.Context, userId int64, reportType defs.ReportType, reportGetFunction reportGetFunc) (string, error) {
	cachedData, err := c.cache.Get(userId, reportType)
	if errors.Is(err, redis.Nil) {
		report, err := reportGetFunction(ctx, userId)
		if err != nil {
			return "", err
		}
		startTime := c.StartTimeByReport(reportType)
		earliestDate, err := c.earliestDateSince(ctx, startTime)
		if err != nil {
			return "", err
		}
		periodDuration := time.Since(c.StartTimeByReport(reportType))
		err = c.cache.Set(userId, reportType, report, earliestDate.Add(periodDuration))
		if err != nil {
			return "", err
		}
		return report, nil
	}
	if err != nil {
		c.logger.Error("error getting report from cache", zap.Error(err))
		return "", err
	}
	return cachedData, nil
}

func (c *cachedReportManager) LastWeek(ctx context.Context, userId int64) (string, error) {
	return c.lastPeriod(ctx, userId, defs.ReportLastWeek, c.rm.LastWeek)
}

func (c *cachedReportManager) LastMonth(ctx context.Context, userId int64) (string, error) {
	return c.lastPeriod(ctx, userId, defs.ReportLastMonth, c.rm.LastMonth)
}

func (c *cachedReportManager) LastYear(ctx context.Context, userId int64) (string, error) {
	return c.lastPeriod(ctx, userId, defs.ReportLastYear, c.rm.LastYear)
}
