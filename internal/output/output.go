package output

type Output interface {
	Currency() CurrencyListOutput
	Reports() ReportManager
	Limits() LimitListOutput
}

type output struct {
	currency CurrencyListOutput
	reports  ReportManager
	limits   LimitListOutput
}

func NewOutput(currency CurrencyListOutput, reports ReportManager, limits LimitListOutput) Output {
	return &output{
		currency: currency,
		reports:  reports,
		limits:   limits,
	}
}

func (o *output) Currency() CurrencyListOutput {
	return o.currency
}

func (o *output) Reports() ReportManager {
	return o.reports
}

func (o *output) Limits() LimitListOutput {
	return o.limits
}
