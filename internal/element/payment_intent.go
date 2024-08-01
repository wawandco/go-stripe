package element

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"stripe-cop/internal/customer"
	"stripe-cop/internal/model"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
)

func HandleCreatePaymentIntent(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_SC_KEY")
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	info := model.PaymentInfo{
		Amount:     "109.99",
		CardHolder: "Edwin Polo",
		Email:      "edwin@example.com",

		BillingLine:    "Theo Parker 123 Pike ST",
		BillingCity:    "Seatle",
		BillingState:   "WA",
		BillingZip:     "98122",
		BillingCountry: "United States",
	}

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(info.AmountInCents())),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		// In the latest version of the API, specifying the `automatic_payment_methods` parameter is optional because Stripe enables its functionality by default.
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Description: stripe.String("One-time payment example, using Stripe element"),
	}

	c, err := customer.CreateCustomer(info)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return
	}
	params.Customer = stripe.String(c)

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return
	}

	writeJSON(w, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	})
}

func HandleCreatePaymentIntentAppFee(w http.ResponseWriter, r *http.Request) {
	stripe.Key = os.Getenv("STRIPE_SC_KEY")
	if r.Method != "POST" {
		http.Error(w, http.StatusText(http.StatusMethodNotAllowed), http.StatusMethodNotAllowed)
		return
	}

	cents := 100
	amount := 100 * cents

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(amount)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		// In the latest version of the API, specifying the `automatic_payment_methods` parameter is optional because Stripe enables its functionality by default.
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		// ApplicationFeeAmount: stripe.Int64(int64(20 * 100)),
		TransferData: &stripe.PaymentIntentTransferDataParams{
			Amount:      stripe.Int64(20 * 100),
			Destination: stripe.String("acct_1PeyvdDK9l7Rmc4L"),
		},
		Description: stripe.String("One-time payment example, using Stripe element. Charging app fee"),
	}

	// params.SetStripeAccount("acct_1PeyvdDK9l7Rmc4L")

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("pi.New: %v", err)
		return
	}

	writeJSON(w, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: pi.ClientSecret,
	})
}

func writeJSON(w http.ResponseWriter, v interface{}) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(v); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		log.Printf("json.NewEncoder.Encode: %v", err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	if _, err := io.Copy(w, &buf); err != nil {
		log.Printf("io.Copy: %v", err)
		return
	}
}
