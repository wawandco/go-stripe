package home

import (
	"fmt"

	"github.com/stripe/stripe-go/v79"
	"github.com/stripe/stripe-go/v79/customer"
)

func CreateCustomer(info PaymentInfo) (string, error) {
	stripe.Key = "sk_test_51IIiV0C5e5WNMZdtXdXmjSkCoEzg1CrCZlweUxjQVGGGDHlGENmCUg1NDhsTgGvgKojTyjVpZXQ2ea6Kk4CCA1to00XQkiBGLq"

	params := &stripe.CustomerParams{
		Name:  stripe.String(info.CardHolder),
		Email: stripe.String(info.Email),
		Address: &stripe.AddressParams{
			City:       stripe.String(info.BillingCity),
			Country:    stripe.String(info.BillingCountry),
			Line1:      stripe.String(info.BillingLine),
			PostalCode: stripe.String(info.BillingZip),
			State:      stripe.String(info.BillingState),
		},
	}

	params.PaymentMethod = stripe.String("pm_card_visa")

	c, err := customer.New(params)
	if err != nil {
		fmt.Printf("Error creating customer: %v", err)
		return "", err
	}

	fmt.Printf("Customer created: %v\n", c.ID)
	return c.ID, nil
}
