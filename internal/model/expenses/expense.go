package expenses

import "time"

type Expense struct {
	UserID   int64
	Category string
	Amount   int
	Currency string
	Date     time.Time
}
