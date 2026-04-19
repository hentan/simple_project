package handlers

import (
	"context"
	"fmt"
)

func (app *Application) GetBalance(ctx context.Context) (int, error) {
	sumExpenses, err := app.DB.GetSumExpenses(ctx)
	if err != nil {
		return 0, fmt.Errorf("get sum expenses: %w", err)
	}

	sumPayments, err := app.DB.GetSumPayments(ctx)
	if err != nil {
		return 0, fmt.Errorf("get sum payments: %w", err)
	}

	summBalance := sumPayments - sumExpenses
	return summBalance, nil
}
