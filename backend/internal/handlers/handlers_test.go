package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"simple_project/internal/config"
	"simple_project/internal/models"
)

// mockDatabase реализует repository.Database полностью.
type mockDatabase struct {
	connectFn func() *sql.DB

	getExpensesFn   func() ([]models.Expense, error)
	addExpenseFn    func(*models.Expense) error
	updateExpenseFn func(*models.Expense) error
	deleteExpenseFn func(*models.Expense) error

	addPaymentFn     func(*models.Payment) error
	updatePaymentFn  func(*models.Payment) error
	deletePaymentFn  func(*models.Payment) error
	getAllPaymentsFn func() ([]models.Payment, error)

	addCalled    int
	updateCalled int
	deleteCalled int
}

func (m *mockDatabase) Connect() *sql.DB {
	if m.connectFn != nil {
		return m.connectFn()
	}
	return nil
}

func (m *mockDatabase) GetExpenses() ([]models.Expense, error) {
	if m.getExpensesFn == nil {
		return nil, nil
	}
	return m.getExpensesFn()
}

func (m *mockDatabase) AddExpense(e *models.Expense) error {
	m.addCalled++
	if m.addExpenseFn == nil {
		return nil
	}
	return m.addExpenseFn(e)
}

func (m *mockDatabase) UpdateExpense(e *models.Expense) error {
	m.updateCalled++
	if m.updateExpenseFn == nil {
		return nil
	}
	return m.updateExpenseFn(e)
}

func (m *mockDatabase) DeleteExpense(e *models.Expense) error {
	m.deleteCalled++
	if m.deleteExpenseFn == nil {
		return nil
	}
	return m.deleteExpenseFn(e)
}

// Методы платежей — для полного соответствия интерфейсу.
func (m *mockDatabase) AddPayment(p *models.Payment) error {
	if m.addPaymentFn == nil {
		return nil
	}
	return m.addPaymentFn(p)
}

func (m *mockDatabase) UpdatePayment(p *models.Payment) error {
	if m.updatePaymentFn == nil {
		return nil
	}
	return m.updatePaymentFn(p)
}

func (m *mockDatabase) DeletePayment(p *models.Payment) error {
	if m.deletePaymentFn == nil {
		return nil
	}
	return m.deletePaymentFn(p)
}

func (m *mockDatabase) GetAllPayments() ([]models.Payment, error) {
	if m.getAllPaymentsFn == nil {
		return nil, nil
	}
	return m.getAllPaymentsFn()
}

func newAppWithMock(db *mockDatabase) *Application {
	return &Application{
		DB:     db,
		config: config.Config{
			// Не важно для unit-тестов хендлеров
		},
	}
}

func TestGetExpenses_OK(t *testing.T) {
	now := time.Date(2026, 1, 26, 10, 0, 0, 0, time.UTC)

	db := &mockDatabase{
		getExpensesFn: func() ([]models.Expense, error) {
			return []models.Expense{
				{
					Id:      1,
					Date:    now,
					GiftFor: "teacher",
					PupilId: 10,
					Summ:    250,
					Surname: "Ivanov",
				},
			}, nil
		},
	}
	app := newAppWithMock(db)

	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	rr := httptest.NewRecorder()

	app.GetExpenses(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))

	var got []models.Expense
	require.NoError(t, json.Unmarshal(rr.Body.Bytes(), &got))
	require.Len(t, got, 1)
	assert.Equal(t, 1, got[0].Id)
	assert.Equal(t, "teacher", got[0].GiftFor)
	assert.Equal(t, 10, got[0].PupilId)
	assert.Equal(t, 250, got[0].Summ)
	assert.Equal(t, "Ivanov", got[0].Surname)
}

func TestGetExpenses_DBError(t *testing.T) {
	db := &mockDatabase{
		getExpensesFn: func() ([]models.Expense, error) {
			return nil, errors.New("db down")
		},
	}
	app := newAppWithMock(db)

	req := httptest.NewRequest(http.MethodGet, "/expenses", nil)
	rr := httptest.NewRecorder()

	app.GetExpenses(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "db down")
}

func TestAddExpense_BadJSON(t *testing.T) {
	db := &mockDatabase{}
	app := newAppWithMock(db)

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBufferString("{bad json"))
	rr := httptest.NewRecorder()

	app.AddExpense(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, 0, db.addCalled)
}

func TestAddExpense_OK(t *testing.T) {
	db := &mockDatabase{
		addExpenseFn: func(e *models.Expense) error {
			// Проверим, что распарсились ключевые поля
			assert.Equal(t, 1, e.Id)
			assert.Equal(t, "mom", e.GiftFor)
			assert.Equal(t, 7, e.PupilId)
			assert.Equal(t, 999, e.Summ)
			assert.Equal(t, "Petrov", e.Surname)
			return nil
		},
	}
	app := newAppWithMock(db)

	// ВАЖНО: time.Time в JSON должен быть RFC3339, иначе decode упадёт.
	body := `{
		"id": 1,
		"date": "2026-01-26T10:00:00Z",
		"gift_for": "mom",
		"pupil_id": 7,
		"summ": 999,
		"surname": "Petrov"
	}`

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	app.AddExpense(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	require.Equal(t, 1, db.addCalled)
}

func TestAddExpense_DBError(t *testing.T) {
	db := &mockDatabase{
		addExpenseFn: func(e *models.Expense) error {
			return errors.New("insert failed")
		},
	}
	app := newAppWithMock(db)

	body := `{
		"id": 1,
		"date": "2026-01-26T10:00:00Z",
		"gift_for": "dad",
		"pupil_id": 7,
		"summ": 100,
		"surname": "Sidorov"
	}`

	req := httptest.NewRequest(http.MethodPost, "/expenses", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	app.AddExpense(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "insert failed")
	require.Equal(t, 1, db.addCalled)
}

func TestUpdateExpense_BadJSON(t *testing.T) {
	db := &mockDatabase{}
	app := newAppWithMock(db)

	req := httptest.NewRequest(http.MethodPut, "/expenses", bytes.NewBufferString("{bad json"))
	rr := httptest.NewRecorder()

	app.UpdateExpense(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, 0, db.updateCalled)
}

func TestUpdateExpense_OK(t *testing.T) {
	db := &mockDatabase{
		updateExpenseFn: func(e *models.Expense) error {
			assert.Equal(t, 2, e.Id)
			assert.Equal(t, 555, e.Summ)
			return nil
		},
	}
	app := newAppWithMock(db)

	body := `{
		"id": 2,
		"date": "2026-01-26T12:00:00Z",
		"gift_for": "friend",
		"pupil_id": 9,
		"summ": 555,
		"surname": "Smirnov"
	}`

	req := httptest.NewRequest(http.MethodPut, "/expenses", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	app.UpdateExpense(rr, req)

	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, "application/json", rr.Header().Get("Content-Type"))
	require.Equal(t, 1, db.updateCalled)
}

func TestUpdateExpense_DBError(t *testing.T) {
	db := &mockDatabase{
		updateExpenseFn: func(e *models.Expense) error {
			return errors.New("update failed")
		},
	}
	app := newAppWithMock(db)

	body := `{
		"id": 2,
		"date": "2026-01-26T12:00:00Z",
		"gift_for": "friend",
		"pupil_id": 9,
		"summ": 555,
		"surname": "Smirnov"
	}`

	req := httptest.NewRequest(http.MethodPut, "/expenses", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	app.UpdateExpense(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "update failed")
	require.Equal(t, 1, db.updateCalled)
}

func TestDeleteExpense_BadJSON(t *testing.T) {
	db := &mockDatabase{}
	app := newAppWithMock(db)

	req := httptest.NewRequest(http.MethodDelete, "/expenses", bytes.NewBufferString("{bad json"))
	rr := httptest.NewRecorder()

	app.DeleteExpense(rr, req)

	require.Equal(t, http.StatusBadRequest, rr.Code)
	require.Equal(t, 0, db.deleteCalled)
}

func TestDeleteExpense_OK_DefaultStatus200(t *testing.T) {
	db := &mockDatabase{
		deleteExpenseFn: func(e *models.Expense) error {
			assert.Equal(t, 3, e.Id)
			return nil
		},
	}
	app := newAppWithMock(db)

	body := `{
		"id": 3,
		"date": "2026-01-26T12:00:00Z",
		"gift_for": "",
		"pupil_id": 0,
		"summ": 0,
		"surname": ""
	}`

	req := httptest.NewRequest(http.MethodDelete, "/expenses", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	app.DeleteExpense(rr, req)

	// В хендлере нет WriteHeader на успешном кейсе, поэтому будет 200 по умолчанию.
	require.Equal(t, http.StatusOK, rr.Code)
	require.Equal(t, 1, db.deleteCalled)
}

func TestDeleteExpense_DBError(t *testing.T) {
	db := &mockDatabase{
		deleteExpenseFn: func(e *models.Expense) error {
			return errors.New("delete failed")
		},
	}
	app := newAppWithMock(db)

	body := `{
		"id": 3,
		"date": "2026-01-26T12:00:00Z",
		"gift_for": "",
		"pupil_id": 0,
		"summ": 0,
		"surname": ""
	}`

	req := httptest.NewRequest(http.MethodDelete, "/expenses", bytes.NewBufferString(body))
	rr := httptest.NewRecorder()

	app.DeleteExpense(rr, req)

	require.Equal(t, http.StatusInternalServerError, rr.Code)
	assert.Contains(t, rr.Body.String(), "delete failed")
	require.Equal(t, 1, db.deleteCalled)
}
