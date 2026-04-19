package handlers

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"testing"

	"simple_project/internal/models"
	"simple_project/internal/repository"

	"github.com/stretchr/testify/require"
)

var errDB = errors.New("db failed")

type fakeDB struct {
	expenses []models.Expense
	payments []models.Payment
	pupils   []models.Pupil

	sumExpenses int
	sumPayments int

	getExpensesErr    error
	addExpenseErr     error
	updateExpenseErr  error
	deleteExpenseErr  error
	getPaymentsErr    error
	addPaymentErr     error
	updatePaymentErr  error
	deletePaymentErr  error
	getSumExpensesErr error
	getSumPaymentsErr error
	getPupilsErr      error
	addPupilErr       error
	updatePupilErr    error
	deletePupilErr    error

	addedExpense   *models.Expense
	updatedExpense *models.Expense
	deletedExpense *models.Expense
	addedPayment   *models.Payment
	updatedPayment *models.Payment
	deletedPayment *models.Payment
	addedPupil     *models.Pupil
	updatedPupil   *models.Pupil
	deletedPupil   *models.Pupil
}

var _ repository.Database = (*fakeDB)(nil)

func (db *fakeDB) GetExpenses(ctx context.Context) ([]models.Expense, error) {
	return db.expenses, db.getExpensesErr
}

func (db *fakeDB) AddExpense(ctx context.Context, expense *models.Expense) error {
	copy := *expense
	db.addedExpense = &copy
	return db.addExpenseErr
}

func (db *fakeDB) UpdateExpense(ctx context.Context, expense *models.Expense) error {
	copy := *expense
	db.updatedExpense = &copy
	return db.updateExpenseErr
}

func (db *fakeDB) DeleteExpense(ctx context.Context, expense *models.Expense) error {
	copy := *expense
	db.deletedExpense = &copy
	return db.deleteExpenseErr
}

func (db *fakeDB) AddPayment(ctx context.Context, payment *models.Payment) error {
	copy := *payment
	db.addedPayment = &copy
	return db.addPaymentErr
}

func (db *fakeDB) UpdatePayment(ctx context.Context, payment *models.Payment) error {
	copy := *payment
	db.updatedPayment = &copy
	return db.updatePaymentErr
}

func (db *fakeDB) DeletePayment(ctx context.Context, payment *models.Payment) error {
	copy := *payment
	db.deletedPayment = &copy
	return db.deletePaymentErr
}

func (db *fakeDB) GetAllPayments(ctx context.Context) ([]models.Payment, error) {
	return db.payments, db.getPaymentsErr
}

func (db *fakeDB) GetSumPayments(ctx context.Context) (int, error) {
	return db.sumPayments, db.getSumPaymentsErr
}

func (db *fakeDB) GetSumExpenses(ctx context.Context) (int, error) {
	return db.sumExpenses, db.getSumExpensesErr
}

func (db *fakeDB) GetPupils(ctx context.Context) ([]models.Pupil, error) {
	return db.pupils, db.getPupilsErr
}

func (db *fakeDB) AddPupil(ctx context.Context, pupil *models.Pupil) error {
	copy := *pupil
	db.addedPupil = &copy
	return db.addPupilErr
}

func (db *fakeDB) UpdatePupil(ctx context.Context, pupil *models.Pupil) error {
	copy := *pupil
	db.updatedPupil = &copy
	return db.updatePupilErr
}

func (db *fakeDB) DeletePupil(ctx context.Context, pupil *models.Pupil) error {
	copy := *pupil
	db.deletedPupil = &copy
	return db.deletePupilErr
}

type spyHandler struct {
	last string
}

func (h *spyHandler) Start(ctx context.Context, handler http.Handler) error {
	return nil
}

func (h *spyHandler) GetExpenses(w http.ResponseWriter, r *http.Request) {
	h.write(w, "GetExpenses")
}

func (h *spyHandler) AddExpense(w http.ResponseWriter, r *http.Request) {
	h.write(w, "AddExpense")
}

func (h *spyHandler) UpdateExpense(w http.ResponseWriter, r *http.Request) {
	h.write(w, "UpdateExpense")
}

func (h *spyHandler) DeleteExpense(w http.ResponseWriter, r *http.Request) {
	h.write(w, "DeleteExpense")
}

func (h *spyHandler) GetPayments(w http.ResponseWriter, r *http.Request) {
	h.write(w, "GetPayments")
}

func (h *spyHandler) AddPayment(w http.ResponseWriter, r *http.Request) {
	h.write(w, "AddPayment")
}

func (h *spyHandler) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	h.write(w, "UpdatePayment")
}

func (h *spyHandler) DeletePayment(w http.ResponseWriter, r *http.Request) {
	h.write(w, "DeletePayment")
}

func (h *spyHandler) GetPupils(w http.ResponseWriter, r *http.Request) {
	h.write(w, "GetPupils")
}

func (h *spyHandler) AddPupil(w http.ResponseWriter, r *http.Request) {
	h.write(w, "AddPupil")
}

func (h *spyHandler) UpdatePupil(w http.ResponseWriter, r *http.Request) {
	h.write(w, "UpdatePupil")
}

func (h *spyHandler) DeletePupil(w http.ResponseWriter, r *http.Request) {
	h.write(w, "DeletePupil")
}

func (h *spyHandler) write(w http.ResponseWriter, name string) {
	h.last = name
	w.WriteHeader(http.StatusOK)
}

func mustJSON(t *testing.T, value any) string {
	t.Helper()

	body, err := json.Marshal(value)
	require.NoError(t, err)
	return string(body)
}
