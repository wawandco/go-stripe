package home

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/leapkit/leapkit/core/render"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/paymentintent"
)

func PayIntentConfirm(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	amount, _ := strconv.Atoi(r.FormValue("amount"))

	info := PaymentInfo{
		Amount:     amount,
		CardHolder: "Javier Hernandez",
		CardNumber: r.FormValue("cnumber"),
		ExpMonth:   r.FormValue("month"),
		ExpYear:    r.FormValue("year"),
		CVC:        r.FormValue("cvc"),

		Email: "javi@example.com",

		BillingLine:    "Theo Parker 123 Pike ST",
		BillingCity:    "Seatle",
		BillingState:   "WA",
		BillingZip:     "98122",
		BillingCountry: "United States",
	}

	// pi, _ := PaymentIntent(info)
	// success, err := ConfirmPaymentIntent(pi)

	pi, err := DirectPaymentIntent(info)
	success := pi != ""

	if err != nil {
		stripeError := StripeError{}
		message := ""
		if jerr := json.Unmarshal([]byte(err.Error()), &stripeError); jerr != nil {
			message = "please review your payment information"
		} else {
			message = stripeError.Message
		}

		rw.Set("error", message)
	}

	rw.Set("success", success)
	rw.Set("target", "example-two")
	rw.Set("backurl", "back-ex-2")
	rw.RenderClean("home/success.html")
}

func BackTwo(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	rw.Set("target", "example-two")
	rw.Set("backurl", "back-ex-2")
	rw.RenderClean("home/s_charge_two.html")
}

func PaymentIntent(info PaymentInfo) (string, error) {
	stripe.Key = os.Getenv("STRIPE_SC_KEY")

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(info.Amount * 100)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		// In the latest version of the API, specifying the `automatic_payment_methods` parameter is optional because Stripe enables its functionality by default.
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled: stripe.Bool(true),
		},
		Confirm:       stripe.Bool(false), // To create the charge
		PaymentMethod: stripe.String("pm_card_visa"),
		Description:   stripe.String("One-time payment example, payment intent confirmed"),
		Metadata:      map[string]string{"Name": "Gopher Toy", "Description": "Toy"},
	}

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		log.Printf("pi.New: %v", err)
		return "", err
	}

	return pi.ID, nil
}

func ConfirmPaymentIntent(piID string) (bool, error) {
	stripe.Key = os.Getenv("STRIPE_SC_KEY")

	params := &stripe.PaymentIntentConfirmParams{
		PaymentMethod: stripe.String("pm_card_visa"),
		ReturnURL:     stripe.String("https://www.example.com"),
	}

	_, err := paymentintent.Confirm(piID, params)
	if err != nil {
		log.Printf("pic.New: %v", err)
		return false, err
	}

	return true, nil
}

func DirectPaymentIntent(info PaymentInfo) (string, error) {
	stripe.Key = os.Getenv("STRIPE_SC_KEY")

	// Create a PaymentIntent with amount and currency
	params := &stripe.PaymentIntentParams{
		Amount:   stripe.Int64(int64(info.Amount * 100)),
		Currency: stripe.String(string(stripe.CurrencyUSD)),
		// In the latest version of the API, specifying the `automatic_payment_methods` parameter is optional because Stripe enables its functionality by default.
		AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			Enabled:        stripe.Bool(true),
			AllowRedirects: stripe.String("never"),
		},
		Confirm:       stripe.Bool(true), // To create the charge
		PaymentMethod: stripe.String("pm_card_visa"),
		Description:   stripe.String("One-time payment example, direct payment intent confirmed"),
		Metadata:      map[string]string{"Name": "Gopher Toy", "Description": "Toy"},
	}

	c, err := CreateCustomer(info)
	if err != nil {
		log.Printf("c.New: %v", err)
		return "", err
	}
	params.Customer = stripe.String(c)

	pi, err := paymentintent.New(params)
	log.Printf("pi.New: %v", pi.ClientSecret)

	if err != nil {
		log.Printf("pi.New: %v", err)
		return "", err
	}

	return pi.ID, nil
}
