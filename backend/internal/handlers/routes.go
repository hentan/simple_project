package handlers

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func Routes(h Handler) http.Handler {
	mux := chi.NewRouter()
	mux.Use(middleware.Recoverer)
	mux.Use(middleware.StripSlashes)
	mux.Get("/expenses", h.GetExpenses)
	mux.Post("/expenses", h.AddExpense)
	mux.Put("/expenses", h.UpdateExpense)
	mux.Delete("/expenses", h.DeleteExpense)
	mux.Get("/payments", h.GetPayments)
	mux.Post("/payments", h.AddPayment)
	mux.Put("/payments", h.UpdatePayment)
	mux.Delete("/payments", h.DeletePayment)
	return mux
}
