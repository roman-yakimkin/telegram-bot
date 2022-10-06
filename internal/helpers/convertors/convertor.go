package convertors

type CurrencyConvertorFrom interface {
	From(amount int, currFrom string) (int, error)
}

type CurrencyConvertorTo interface {
	To(amount int, currTo string) (int, error)
}

type CurrencyConvertor interface {
	CurrencyConvertorFrom
	CurrencyConvertorTo
}
