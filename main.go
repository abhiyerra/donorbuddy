package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"

	"github.com/stripe/stripe-go"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"

	"gopkg.in/pg.v5"
)

var config struct {
	Database *pg.Options

	Auth struct {
		SecurityKey string

		Facebook struct {
			AppID     string
			AppSecret string
			Callback  string
		}
	}

	StripeSecretKey string
	// StripePlan should be a Plan in Stripe with the price of
	// 0.01 and billed monthly.
	StripePlan string

	DB *pg.DB `json:"-"`
}

func readConfig() {
	configFile, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal("failed to open config file", os.Args[1])
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatal("failed to open config file", os.Args[1])
	}

	log.Println(config)
}

func setConfig() {
	if config.Database == nil {
		log.Fatal("failed to load database configuration")
	}

	config.DB = pg.Connect(config.Database)

	stripe.Key = config.StripeSecretKey

	gomniauth.SetSecurityKey(config.Auth.SecurityKey)
	gomniauth.WithProviders(facebook.New(config.Auth.Facebook.AppID, config.Auth.Facebook.AppSecret, config.Auth.Facebook.Callback))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
}

func main() {
	readConfig()
	setConfig()

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

	http.Handle("/", handlers.LoggingHandler(os.Stdout, r))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
