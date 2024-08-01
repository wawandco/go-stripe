package model

import (
	"math"
	"strconv"
)

type StripeError struct {
	Code    string `json:"code"`
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type PaymentInfo struct {
	Email string

	// card info
	Amount     string
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

func (pi PaymentInfo) AmountInCents() int {
	amount, err := strconv.ParseFloat(pi.Amount, 64)
	if err != nil {
		return 0
	}

	cents := math.Round(amount * 100)
	return int(cents)
}
