package localerr

import "github.com/pkg/errors"

var ErrUserStateNotFound = errors.New("user state not found")

var ErrCurrencyNotFound = errors.New("currency not found")
var ErrCannotGetCurrencyRates = errors.New("cannot get currency rates")

var ErrIncorrectAmountValue = errors.New("incorrect amount value")
