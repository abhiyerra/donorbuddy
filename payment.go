package main

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

	s, err := sub.New(&stripe.SubParams{
		Customer: "cus_9sek9eRTNJ0BdG",
		Plan:     config.StripePlan,
		Quantity: plan.Amount,
	})

}

func updatePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	var plan struct {
		// NewAmount is the amount User wants to donate in cents
		NewAmount uint64
	}

	// TODO: Check amount > 0

	s, err := sub.Update(
		"sub_9sed4J2K4jurwS", // TODO: Get User's subscription
		&stripe.SubParams{
			Plan:     config.StripePlan,
			Quantity: plan.NewAmount,
		},
	)

	respondJson(w, r, struct{}{})
}

func deletePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	err := sub.Cancel(
		"sub_9sed4J2K4jurwS", // TODO: Cancel user's subscription
	)

	respondJson(w, r, struct{}{})
}
