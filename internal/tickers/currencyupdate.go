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
	earlyDate time.Time
}

func NewCurrencyUpdate(cr repo.CurrencyRateRepo, seconds time.Duration, earlyDate time.Time) *CurrencyUpdate {
	return &CurrencyUpdate{
		cr:        cr,
		seconds:   seconds,
		earlyDate: earlyDate,
	}
}

func (c *CurrencyUpdate) Run(ctx context.Context) {
	ticker := time.NewTicker(c.seconds)
	date := utils.TimeTruncate(time.Now())
	go func(ctx context.Context, t *time.Ticker) {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				hasData := true
				for hasData && date.After(c.earlyDate) {
					curr, err := c.cr.GetAllByDate(date)
					if err != nil {
						log.Print("Error upon getting currency rates:", err)
						time.Sleep(5 * time.Second)
					}
					hasData = len(curr) > 0
					if hasData {
						date = date.AddDate(0, 0, -1)
					}
				}
				if !date.After(c.earlyDate) {
					err := c.cr.LoadByDate(date)
					if err != nil {
						log.Print("Error upon getting currency rates:", err)
						time.Sleep(5 * time.Second)
					}
				}
			}
		}
	}(ctx, ticker)
}
