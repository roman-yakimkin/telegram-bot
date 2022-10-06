package output

type Output struct {
	currency *CurrencyListOutput
	reports  *ReportManager
}

func NewOutput(currency *CurrencyListOutput, reports *ReportManager) *Output {
	return &Output{
		currency: currency,
		reports:  reports,
	}
}

func (o *Output) Currency() *CurrencyListOutput {
	return o.currency
}

func (o *Output) Reports() *ReportManager {
	return o.reports
}
