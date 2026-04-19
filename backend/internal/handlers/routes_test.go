package handlers

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRoutes(t *testing.T) {
	handler := &spyHandler{}
	router := Routes(handler)

	tests := []struct {
		name   string
		method string
		path   string
		want   string
	}{
		{name: "get expenses", method: http.MethodGet, path: "/expenses", want: "GetExpenses"},
		{name: "post expenses", method: http.MethodPost, path: "/expenses", want: "AddExpense"},
		{name: "put expenses", method: http.MethodPut, path: "/expenses", want: "UpdateExpense"},
		{name: "delete expenses", method: http.MethodDelete, path: "/expenses", want: "DeleteExpense"},
		{name: "get payments", method: http.MethodGet, path: "/payments", want: "GetPayments"},
		{name: "post payments", method: http.MethodPost, path: "/payments", want: "AddPayment"},
		{name: "put payments", method: http.MethodPut, path: "/payments", want: "UpdatePayment"},
		{name: "delete payments", method: http.MethodDelete, path: "/payments", want: "DeletePayment"},
		{name: "get pupils", method: http.MethodGet, path: "/pupils", want: "GetPupils"},
		{name: "post pupils", method: http.MethodPost, path: "/pupils", want: "AddPupil"},
		{name: "put pupils", method: http.MethodPut, path: "/pupils", want: "UpdatePupil"},
		{name: "delete pupils", method: http.MethodDelete, path: "/pupils", want: "DeletePupil"},
		{name: "strip trailing slash", method: http.MethodGet, path: "/expenses/", want: "GetExpenses"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			handler.last = ""
			req := httptest.NewRequest(tt.method, tt.path, nil)
			rr := httptest.NewRecorder()

			router.ServeHTTP(rr, req)

			require.Equal(t, http.StatusOK, rr.Code)
			require.Equal(t, tt.want, handler.last)
		})
	}
}
