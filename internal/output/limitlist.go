package output

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gitlab.ozon.dev/r.yakimkin/telegram-bot/internal/repo"
)

type LimitListOutput interface {
	Output(ctx context.Context, UserID int64) (string, error)
}

type limitListOutput struct {
	limitRepo    repo.ExpenseLimitsRepo
	outputAmount CurrencyAmount
}

func NewLimitListOutput(limitRepo repo.ExpenseLimitsRepo, outputAmount CurrencyAmount) LimitListOutput {
	return &limitListOutput{
		limitRepo:    limitRepo,
		outputAmount: outputAmount,
	}
}

func (o *limitListOutput) Output(ctx context.Context, UserID int64) (string, error) {
	var sb strings.Builder
	limits, err := o.limitRepo.GetAll(ctx, UserID)
	if err != nil {
		return "", err
	}
	for i, limit := range limits {
		amount, err := o.outputAmount.Output(ctx, limit.Value, "RUB")
		if err != nil {
			return "", err
		}
		sb.WriteString(fmt.Sprintf("%02d (%s) - %s\n", i+1, time.Month(i+1).String(), amount))
	}
	return sb.String(), nil
}
