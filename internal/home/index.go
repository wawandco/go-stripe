package home

import (
	"net/http"

	"github.com/leapkit/leapkit/core/render"
)

func Index(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	err := rw.Render("home/index.html")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

type StripeError struct {
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type PaymentInfo struct {
	Email string

	// card info
	Amount     int
	CardHolder string
	CardNumber string
	CVC        string
	ExpMonth   string
	ExpYear    string

	// billing info
	BillingLine    string
	BillingCity    string
	BillingState   string
	BillingZip     string
	BillingCountry string
}
