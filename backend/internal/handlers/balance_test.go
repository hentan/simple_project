package handlers

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestApplicationGetBalance(t *testing.T) {
	tests := []struct {
		name            string
		db              *fakeDB
		expectedBalance int
		requireErr      require.ErrorAssertionFunc
	}{
		{
			name:            "subtracts expenses from payments",
			db:              &fakeDB{sumPayments: 1000, sumExpenses: 350},
			expectedBalance: 650,
			requireErr:      require.NoError,
		},
		{
			name:            "allows negative balance",
			db:              &fakeDB{sumPayments: 100, sumExpenses: 250},
			expectedBalance: -150,
			requireErr:      require.NoError,
		},
		{
			name:       "returns expense sum error",
			db:         &fakeDB{getSumExpensesErr: errDB},
			requireErr: require.Error,
		},
		{
			name:       "returns payment sum error",
			db:         &fakeDB{getSumPaymentsErr: errDB},
			requireErr: require.Error,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{DB: tt.db}

			balance, err := app.GetBalance(context.Background())

			tt.requireErr(t, err)
			require.Equal(t, tt.expectedBalance, balance)
		})
	}
}
