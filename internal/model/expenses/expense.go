package expenses

import "time"

type Expense struct {
	UserId   int64
	Category string
	Amount   int
	Currency string
	Date     time.Time
}
