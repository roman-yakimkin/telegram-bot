package convertors

import (
	"context"
	"time"
)

type CurrencyConvertorFrom interface {
	From(ctx context.Context, amount int, currFrom string, date time.Time) (int, error)
}

type CurrencyConvertorTo interface {
	To(ctx context.Context, amount int, currTo string, date time.Time) (int, error)
}

type CurrencyConvertor interface {
	CurrencyConvertorFrom
	CurrencyConvertorTo
}
