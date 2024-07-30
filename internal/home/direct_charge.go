package home

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"

	"github.com/leapkit/leapkit/core/render"
	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/charge"
)

func PayChargeOne(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	amount, _ := strconv.Atoi(r.FormValue("amount"))

	info := PaymentInfo{
		Amount:     amount,
		CardHolder: r.FormValue("cardholder"),
		CardNumber: r.FormValue("cnumber"),
		ExpMonth:   r.FormValue("month"),
		ExpYear:    r.FormValue("year"),
		CVC:        r.FormValue("cvc"),

		Email: r.FormValue("email"),

		BillingLine:    "Theo Parker 123 Pike ST",
		BillingCity:    "Seatle",
		BillingState:   "WA",
		BillingZip:     "98122",
		BillingCountry: "United States",
	}

	success, err := PaymentCharge(info)
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

	rw.Set("info", info)
	rw.Set("success", success)
	rw.Set("target", "example-one")
	rw.Set("backurl", "back-ex-1")
	rw.RenderClean("home/success.html")
}

func BackOne(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	rw.Set("target", "example-one")
	rw.Set("backurl", "back-ex-1")
	rw.RenderClean("home/s_charge_one.html")
}

func PaymentCharge(info PaymentInfo) (bool, error) {
	stripe.Key = os.Getenv("STRIPE_SC_KEY")

	cents := 100
	amount := info.Amount * cents

	chargeParams := &stripe.ChargeParams{
		Amount:      stripe.Int64(int64(amount)),
		Currency:    stripe.String(string(stripe.CurrencyUSD)),
		Description: stripe.String("One-time payment example, direct charge"),
		Capture:     stripe.Bool(true),
		Metadata:    map[string]string{"Name": "Gopher Toy", "Description": "Toy"},
	}

	/*
		//This is the older way to set credit card info to a charge.
		t := &stripe.TokenParams{
			Card: &stripe.CardParams{
				Name:           stripe.String(info.CardHolder),
				Number:         stripe.String(info.CardNumber),
				ExpMonth:       stripe.String(info.ExpMonth),
				ExpYear:        stripe.String(info.ExpYear),
				CVC:            stripe.String(info.CVC),
				AddressLine1:   stripe.String(info.BillingLine),
				AddressCity:    stripe.String(info.BillingCity),
				AddressState:   stripe.String(info.BillingState),
				AddressZip:     stripe.String(info.BillingZip),
				AddressCountry: stripe.String(info.BillingCountry),
			},
		}

		nt, err := token.New(t)
		if err != nil {
			fmt.Printf("Error creating token for charge: %v", err)
			return false, err
		}

		chargeParams.SetSource(nt.ID)
	*/

	chargeParams.SetSource("tok_visa")

	ch, err := charge.New(chargeParams)
	if err != nil {
		fmt.Printf("Error creating charge: %v", err)
		return false, err
	}

	fmt.Printf("Charge created: %v\n", ch.ID)
	return true, nil
}

func PayChargeWithAppFee(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	amount, _ := strconv.Atoi(r.FormValue("amount"))

	info := PaymentInfo{
		Amount:     amount,
		CardHolder: r.FormValue("cardholder"),
		CardNumber: r.FormValue("cnumber"),
		ExpMonth:   r.FormValue("month"),
		ExpYear:    r.FormValue("year"),
		CVC:        r.FormValue("cvc"),

		Email: r.FormValue("email"),

		BillingLine:    "Theo Parker 123 Pike ST",
		BillingCity:    "Seatle",
		BillingState:   "WA",
		BillingZip:     "98122",
		BillingCountry: "United States",
	}

	success, err := PaymentChargeAPPFee(info)
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

	rw.Set("info", info)
	rw.Set("success", success)
	rw.Set("target", "example-four")
	rw.Set("backurl", "back-ex-4")
	rw.RenderClean("home/success.html")
}

func PaymentChargeAPPFee(info PaymentInfo) (bool, error) {
	stripe.Key = os.Getenv("STRIPE_SC_KEY")

	cents := 100
	amount := info.Amount * cents

	chargeParams := &stripe.ChargeParams{
		Amount:               stripe.Int64(int64(amount)),
		Currency:             stripe.String(string(stripe.CurrencyUSD)),
		Description:          stripe.String("One-time payment example, direct charge to connect account with app fee"),
		Capture:              stripe.Bool(true),
		Metadata:             map[string]string{"Name": "Gopher Toy", "Description": "Toy"},
		ApplicationFeeAmount: stripe.Int64(20 * 100),
	}

	chargeParams.SetSource("tok_visa")

	chargeParams.SetStripeAccount("acct_1PeyvdDK9l7Rmc4L")

	ch, err := charge.New(chargeParams)
	if err != nil {
		fmt.Printf("Error creating charge: %v", err)
		return false, err
	}

	fmt.Printf("Charge created: %v\n", ch.ID)
	return true, nil
}

func BackFour(w http.ResponseWriter, r *http.Request) {
	rw := render.FromCtx(r.Context())

	rw.Set("target", "example-four")
	rw.Set("backurl", "back-ex-4")
	rw.RenderClean("home/s_charge_one_app_fee.html")
}
