package handlers

import (
	"context"
	"fmt"
)

func (app *Application) GetBalance(ctx context.Context) (int, error) {
	expensesTotal, err := app.DB.GetExpensesTotal(ctx)
	if err != nil {
		return 0, fmt.Errorf("get expenses total: %w", err)
	}

	paymentsTotal, err := app.DB.GetPaymentsTotal(ctx)
	if err != nil {
		return 0, fmt.Errorf("get payments total: %w", err)
	}

	balance := paymentsTotal - expensesTotal
	return balance, nil
}
