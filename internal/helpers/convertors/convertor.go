package convertors

import "time"

type CurrencyConvertorFrom interface {
	From(amount int, currFrom string, date time.Time) (int, error)
}

type CurrencyConvertorTo interface {
	To(amount int, currTo string, date time.Time) (int, error)
}

type CurrencyConvertor interface {
	CurrencyConvertorFrom
	CurrencyConvertorTo
}
