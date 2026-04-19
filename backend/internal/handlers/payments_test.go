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

func TestApplicationGetPayments(t *testing.T) {
	date := time.Date(2026, 4, 2, 12, 0, 0, 0, time.UTC)

	tests := []struct {
		name       string
		db         *fakeDB
		wantStatus int
		assertBody func(t *testing.T, body *bytes.Buffer)
	}{
		{
			name: "returns payments",
			db: &fakeDB{payments: []models.Payment{{
				ID: 1, Date: date, PupilID: 2, Amount: 500, Surname: "Petrov", Purpose: "monthly",
			}}},
			wantStatus: http.StatusOK,
			assertBody: func(t *testing.T, body *bytes.Buffer) {
				var response []models.Payment
				require.NoError(t, json.NewDecoder(body).Decode(&response))
				require.Len(t, response, 1)
				require.Equal(t, "monthly", response[0].Purpose)
			},
		},
		{
			name:       "returns payments db error",
			db:         &fakeDB{getPaymentsErr: errDB},
			wantStatus: http.StatusInternalServerError,
			assertBody: func(t *testing.T, body *bytes.Buffer) {
				require.Contains(t, body.String(), errDB.Error())
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{DB: tt.db}
			req := httptest.NewRequest(http.MethodGet, "/payments", nil)
			rr := httptest.NewRecorder()

			app.GetPayments(rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			tt.assertBody(t, rr.Body)
		})
	}
}

func TestApplicationPaymentMutations(t *testing.T) {
	payment := models.Payment{ID: 6, PupilID: 7, Amount: 800, Purpose: "class fund"}

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
			name:       "add payment",
			method:     http.MethodPost,
			body:       mustJSON(t, payment),
			db:         &fakeDB{},
			call:       (*Application).AddPayment,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.addedPayment)
				require.Equal(t, payment.Purpose, db.addedPayment.Purpose)
			},
		},
		{
			name:       "update payment",
			method:     http.MethodPut,
			body:       mustJSON(t, payment),
			db:         &fakeDB{},
			call:       (*Application).UpdatePayment,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.updatedPayment)
				require.Equal(t, payment.ID, db.updatedPayment.ID)
			},
		},
		{
			name:       "delete payment",
			method:     http.MethodDelete,
			body:       mustJSON(t, payment),
			db:         &fakeDB{},
			call:       (*Application).DeletePayment,
			wantStatus: http.StatusOK,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.deletedPayment)
				require.Equal(t, payment.ID, db.deletedPayment.ID)
			},
		},
		{
			name:       "bad json does not call db",
			method:     http.MethodPost,
			body:       `{"pupil_id":`,
			db:         &fakeDB{},
			call:       (*Application).AddPayment,
			wantStatus: http.StatusBadRequest,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.Nil(t, db.addedPayment)
			},
		},
		{
			name:       "db error",
			method:     http.MethodPut,
			body:       mustJSON(t, payment),
			db:         &fakeDB{updatePaymentErr: errDB},
			call:       (*Application).UpdatePayment,
			wantStatus: http.StatusInternalServerError,
			assertDB: func(t *testing.T, db *fakeDB) {
				require.NotNil(t, db.updatedPayment)
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			app := &Application{DB: tt.db}
			req := httptest.NewRequest(tt.method, "/payments", bytes.NewBufferString(tt.body))
			rr := httptest.NewRecorder()

			tt.call(app, rr, req)

			require.Equal(t, tt.wantStatus, rr.Code)
			tt.assertDB(t, tt.db)
		})
	}
}
