package handlers

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"simple_project/internal/models"

	"github.com/stretchr/testify/require"
)

func TestApplicationGetExpenses(t *testing.T) {
	date := time.Date(2026, 4, 1, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		db         *fakeDB
		wantStatus int
		assertBody func(t *testing.T, body *bytes.Buffer)
	}{
		{
			name: "returns expenses with balance",
			db: &fakeDB{
				expenses:      []models.Expense{{ID: 1, Date: date, Purpose: "books", PupilID: 2, Amount: 300, Surname: "Ivanov"}},
				paymentsTotal: 1000,
				expensesTotal: 300,
			},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, body *bytes.Buffer) {
				var response models.ExpenseWithBalance
				require.NoError(t, json.NewDecoder(body).Decode(&response))
				require.Equal(t, 700, response.Balance)
				require.Len(t, response.Expenses, 1)
				require.Equal(t, "books", response.Expenses[0].Purpose)
			},
		},
		{
			name:       "returns empty list when repository has no expenses",
			db:         &fakeDB{},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, body *bytes.Buffer) {
				var response models.ExpenseWithBalance
				require.NoError(t, json.NewDecoder(body).Decode(&response))
				require.Empty(t, response.Expenses)
				require.Equal(t, 0, response.Balance)
			},
		},
		{
			name:       "returns db error from expenses query",
			db:         &fakeDB{getExpensesErr: errDB},
			wantStatus: http.StatusInternalServerError,
			assertBody: func(t *testing.T, body *bytes.Buffer) {
				require.Contains(t, body.String(), errDB.Error())
			},
		},
		{
			name:       "returns db error from balance query",
			db:         &fakeDB{getSumPaymentsErr: errDB},
			wantStatus: http.StatusInternalServerError,
			assertBody: func(t *testing.T, body *bytes.Buffer) {
				require.Contains(t, body.String(), "get payments total")
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{DB: tt.db}
			req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
			rr := httptest.NewRecorder()

			app.GetExpenses(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			tt.assertBody(t, rr.Body)
		})
	}
}

func TestApplicationExpenseMutations(t *testing.T) {
	expense := models.Expense{ID: 4, Purpose: "flowers", PupilID: 5, Amount: 250}

	tests := []struct {
		name       string
		method     string
		body       string
		db         *fakeDB
		call       func(app *Application, w http.ResponseWriter, req *http.Request)
		wantStatus int
		assertDB   func(t *testing.T, db *fakeDB)
	}{
		{
			name:       "add expense",
			method:     http.MethodPost,
			body:       mustJSON(t, expense),
			db:         &fakeDB{},
			call:       (*Application).AddExpense,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.addedExpense)
				require.Equal(t, expense.Purpose, db.addedExpense.Purpose)
			},
		},
		{
			name:       "update expense",
			method:     http.MethodPut,
			body:       mustJSON(t, expense),
			db:         &fakeDB{},
			call:       (*Application).UpdateExpense,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.updatedExpense)
				require.Equal(t, expense.ID, db.updatedExpense.ID)
			},
		},
		{
			name:       "delete expense",
			method:     http.MethodDelete,
			body:       mustJSON(t, expense),
			db:         &fakeDB{},
			call:       (*Application).DeleteExpense,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.deletedExpense)
				require.Equal(t, expense.ID, db.deletedExpense.ID)
			},
		},
		{
			name:       "bad json does not call db",
			method:     http.MethodPost,
			body:       `{"id":`,
			db:         &fakeDB{},
			call:       (*Application).AddExpense,
			wantStatus: http.StatusBadRequest,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.Nil(t, db.addedExpense)
			},
		},
		{
			name:       "db error",
			method:     http.MethodPost,
			body:       mustJSON(t, expense),
			db:         &fakeDB{addExpenseErr: errDB},
			call:       (*Application).AddExpense,
			wantStatus: http.StatusInternalServerError,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.addedExpense)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{DB: tt.db}
			req := httptest.NewRequest(tt.method, "/expenses", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			tt.call(app, rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			tt.assertDB(t, tt.db)
		})
	}
}
