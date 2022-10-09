package currencies

import "time"

type CurrencyRate struct {
	Name       string
	RateToMain float64
	Date       time.Time
}
