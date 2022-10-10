package output

type Output interface {
	Currency() CurrencyListOutput
	Reports() ReportManager
}

type output struct {
	currency CurrencyListOutput
	reports  ReportManager
}

func NewOutput(currency CurrencyListOutput, reports ReportManager) Output {
	return &output{
		currency: currency,
		reports:  reports,
	}
}

func (o *output) Currency() CurrencyListOutput {
	return o.currency
}

func (o *output) Reports() ReportManager {
	return o.reports
}
