package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"
)

func createPaymentsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		user = UserValue(r)
	)

	r.ParseForm()
	log.Println(r.Form)

	if user.StripeCustomerID == "" {
		customerParams := &stripe.CustomerParams{
			Email: r.FormValue("email"),
			Desc:  fmt.Sprintf("Donor: %s <%s>", r.FormValue("name"), r.FormValue("email")),
		}
		customerParams.SetSource(r.FormValue("stripeToken"))
		c, err := customer.New(customerParams)
		if err != nil {
			respondJson(w, r, err)
			return
		}

		user.StripeCustomerID = c.ID
		_, err = config.DB.Model(&user).Column("stripe_customer_id").Column("updated_at").Update()
		if err != nil {
			respondJson(w, r, err)
			return
		}
	}

	amount, err := strconv.ParseUint(r.FormValue("amount"), 10, 64)
	if err != nil {
		respondJson(w, r, err)
		return
	}

	s, err := sub.New(&stripe.SubParams{
		Customer: user.StripeCustomerID,
		Plan:     config.StripePlan,
		Quantity: amount,
	})
	if err != nil {
		respondJson(w, r, err)
		return
	}

	user.StripeSubscriptionID = s.ID
	_, err = config.DB.Model(&user).Column("stripe_subscription_id").Column("updated_at").Update()
	if err != nil {
		respondJson(w, r, err)
		return
	}

	http.Redirect(w, r, r.Referer(), http.StatusSeeOther)
}

// func updatePaymentsHandler(w http.ResponseWriter, r *http.Request) {
// 	var plan struct {
// 		// NewAmount is the amount User wants to donate in cents
// 		NewAmount uint64
// 	}

// 	// TODO: Check amount > 0

// 	_, err := sub.Update(
// 		"sub_9sed4J2K4jurwS", // TODO: Get User's subscription
// 		&stripe.SubParams{
// 			Plan:     config.StripePlan,
// 			Quantity: plan.NewAmount,
// 		},
// 	)
// 	if err != nil {
// 		respondJson(w, r, err)
// 		return
// 	}

// 	respondJson(w, r, struct{}{})
// }

func deletePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	_, err := sub.Cancel("sub_9sed4J2K4jurwS", nil) // TODO: Cancel user's subscription
	if err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, struct{}{})
}
