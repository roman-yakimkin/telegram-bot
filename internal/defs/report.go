package defs

type ReportType string

func (r ReportType) String() string {
	return string(r)
}

const (
	ReportLastWeek  ReportType = "lastweek"
	ReportLastMonth ReportType = "lastmonth"
	ReportLastYear  ReportType = "lastyear"
)
