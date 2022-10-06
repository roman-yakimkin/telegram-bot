package tickers

import (
	"context"
	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
	"log"
	"time"
)

type CurrencyUpdate struct {
	cr      repo.CurrencyRepo
	seconds time.Duration
}

func NewCurrencyUpdate(cr repo.CurrencyRepo, seconds time.Duration) *CurrencyUpdate {
	return &CurrencyUpdate{
		cr:      cr,
		seconds: seconds,
	}
}

func (c *CurrencyUpdate) Run(ctx context.Context) {
	ticker := time.NewTicker(c.seconds)
	go func(ctx context.Context, t *time.Ticker) {
		for {
			select {
			case <-ctx.Done():
				t.Stop()
				return
			case <-t.C:
				err := c.cr.LoadAll()
				if err != nil {
					log.Println("loading rates error: ", err)
					t.Stop()
					return
				}
			}
		}
	}(ctx, ticker)
}
