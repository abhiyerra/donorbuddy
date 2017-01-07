package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"

	"github.com/stripe/stripe-go"

	"github.com/gorilla/mux"
	"gopkg.in/pg.v5"
)

var config struct {
	StripeSecretKey string
	// StripePlan should be a Plan in Stripe with the price of
	// 0.01 and billed monthly.
	StripePlan string

	DB *pg.DB
}

func main() {
	config.DB = pg.Connect(&pg.Options{
		User: "postgres",
	})

	stripe.Key = config.StripeSecretKey

	// TODO Fix this
	gomniauth.SetSecurityKey("yLiCQYG7CAflDavqGH461IO0MHp7TEbpg6TwHBWdJzNwYod1i5ZTbrIF5bEoO3oP") // NOTE: DO NOT COPY THIS - MAKE YOR OWN!
	gomniauth.WithProviders(
		// TODO Move this to config and get actual keys.
		facebook.New("537611606322077", "f9f4d77b3d3f4f5775369f5c9f88f65e", "http://localhost:8080/auth/facebook/callback"),
	)

	r := mux.NewRouter()
	r.HandleFunc("/auth/facebook", loginHandler)
	r.HandleFunc("/auth/facebook/callback", loginCallbackHandler)

	r.HandleFunc("/v1/orgs/{orgId}", showOrgsHandler).Methods("GET")
	r.HandleFunc("/v1/orgs", searchOrgsHandler).Methods("GET")

	r.HandleFunc("/v1/payments", createPaymentsHandler).Methods("POST")
	r.HandleFunc("/v1/payments", updatePaymentsHandler).Methods("UPDATE")
	r.HandleFunc("/v1/payments", deletePaymentsHandler).Methods("DELETE")
	//	r.HandleFunc("/v1/payments/stripe-callback", callbackPaymentsHandler).Methods("POST")

	r.HandleFunc("/v1/user/orgs/{orgId}", putUserOrgsHandler).Methods("PUT")
	r.HandleFunc("/v1/user/orgs/{orgId}", deleteUserOrgsHandler).Methods("DELETE")

	r.HandleFunc("/v1/user", showUserHandler)

	http.Handle("/", r)
}
