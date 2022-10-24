package localerr

import "github.com/pkg/errors"

var ErrUserStateNotFound = errors.New("user state not found")

var ErrCurrencyNotFound = errors.New("currency not found")
var ErrCannotGetCurrencyRates = errors.New("cannot get currency rates")
var ErrCurrencyRateUnset = errors.New("currency rate is absent")

var ErrIncorrectAmountValue = errors.New("incorrect amount value")

var ErrExpenseLimitNotFound = errors.New("expense limit not found")

var ErrNotNumericValue = errors.New("not numeric value")
