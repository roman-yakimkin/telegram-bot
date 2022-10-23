package tickers

import (
	"context"
	"log"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/helpers/utils"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type CurrencyUpdate struct {
	cr        repo.CurrencyRateRepo
	seconds   time.Duration
	daysCount int
}

func NewCurrencyUpdate(cr repo.CurrencyRateRepo, seconds time.Duration, daysCount int) *CurrencyUpdate {
	return &CurrencyUpdate{
		cr:        cr,
		seconds:   seconds,
		daysCount: daysCount,
	}
}

func (c *CurrencyUpdate) Run(ctx context.Context) {
	ticker := time.NewTicker(c.seconds)
	date := utils.TimeTruncate(time.Now())
	startDate := date.AddDate(0, 0, -c.daysCount)
	go func(ctx context.Context, t *time.Ticker) {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				hasData := true
				var err error
				for hasData && date.After(startDate) {
					hasData, err = c.cr.HasRatesByDate(ctx, date)
					if err != nil {
						log.Print("Error upon checking rates:", err)
					}
					if hasData {
						date = date.AddDate(0, 0, -1)
					}
				}
				if date.After(startDate) {
					err = c.cr.LoadByDateIfEmpty(ctx, date)
					if err != nil {
						log.Print("Error upon getting currency rates:", err)
					}
				}
			}
		}
	}(ctx, ticker)
}
