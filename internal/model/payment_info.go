package model

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
