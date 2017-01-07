package main

import (
	"net/http"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
)

func createPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	var plan struct {
		// Amount is the amount User wants to donate in cents
		Amount      uint64
		StripeToken string
	}

	customerParams := &stripe.CustomerParams{
		Desc: "Customer for jacob.jackson@example.com",
	}
	customerParams.SetSource(plan.StripeToken) // obtained with Stripe.js
	c, err := customer.New(customerParams)
	if err != nil {
		respondJson(w, r, err)
		return
	}

	// config.DB.Update("") customerid

	s, err := sub.New(&stripe.SubParams{
		Customer: c.ID,
		Plan:     config.StripePlan,
		Quantity: plan.Amount,
	})
	if err != nil {
		respondJson(w, r, err)
		return
	}

	// Update subscriptionid on the user
	_ = s

	respondJson(w, r, struct{}{})
}

func updatePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	var plan struct {
		// NewAmount is the amount User wants to donate in cents
		NewAmount uint64
	}

	// TODO: Check amount > 0

	_, err := sub.Update(
		"sub_9sed4J2K4jurwS", // TODO: Get User's subscription
		&stripe.SubParams{
			Plan:     config.StripePlan,
			Quantity: plan.NewAmount,
		},
	)
	if err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, struct{}{})
}

func deletePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := sub.Cancel("sub_9sed4J2K4jurwS", nil) // TODO: Cancel user's subscription
	if err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, struct{}{})
}
