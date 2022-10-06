package currencies

import "time"

type Currency struct {
	Name       string
	Display    string
	RateToMain float64
	Received   time.Time
}
