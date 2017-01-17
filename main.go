package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"

	"github.com/stripe/stripe-go"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"

	"github.com/rs/cors"

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

	DB        *pg.DB                `json:"-"`
	OAuthConf *oauth2.Config        `json:"-"`
	Store     *sessions.CookieStore `json:"-"`
}

func readConfig() {
	configEnv := os.Getenv("DONORBUDDY_CONFIG")
	if configEnv == "" {
		err := json.Unmarshal([]byte(configEnv), &config)
		if err != nil {
			log.Fatal("failed to open config file", os.Args[1])
		}
	} else {
		configFile, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			log.Fatal("failed to open config file", os.Args[1])
		}

		err = json.Unmarshal(configFile, &config)
		if err != nil {
			log.Fatal("failed to open config file", os.Args[1])
		}
	}
	log.Println(config)
}

func setConfig() {
	if config.Database == nil {
		log.Fatal("failed to load database configuration")
	}

	config.DB = pg.Connect(config.Database)

	stripe.Key = config.StripeSecretKey

	config.OAuthConf = &oauth2.Config{
		ClientID:     config.Auth.Facebook.AppID,
		ClientSecret: config.Auth.Facebook.AppSecret,
		RedirectURL:  config.Auth.Facebook.Callback,
		Scopes:       []string{"public_profile"},
		Endpoint:     facebook.Endpoint,
	}

	config.Store = sessions.NewCookieStore([]byte(config.Auth.SecurityKey))
}

func init() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	pg.SetQueryLogger(log.New(os.Stdout, "", log.LstdFlags))

}

func main() {
	readConfig()
	setConfig()

	r := mux.NewRouter()
	r.HandleFunc("/auth/login", loginHandler)
	r.HandleFunc("/auth/facebook/callback", loginCallbackHandler)

	r.HandleFunc("/v1/orgs/{orgId}", showOrgsHandler).Methods("GET")
	r.HandleFunc("/v1/orgs", searchOrgsHandler).Methods("GET")

	r.Handle("/v1/payments", AuthMiddleware(http.HandlerFunc(createPaymentsHandler))).Methods("POST")
	//r.Handle("/v1/payments", AuthMiddleware(http.HandlerFunc(updatePaymentsHandler))).Methods("UPDATE")
	r.Handle("/v1/payments", AuthMiddleware(http.HandlerFunc(deletePaymentsHandler))).Methods("DELETE")
	//	r.HandleFunc("/v1/payments/stripe-callback", callbackPaymentsHandler).Methods("POST")

	r.Handle("/v1/user/org/{orgId}", AuthMiddleware(http.HandlerFunc(putUserOrgsHandler))).Methods("PUT")
	r.Handle("/v1/user/org/{orgId}", AuthMiddleware(http.HandlerFunc(deleteUserOrgsHandler))).Methods("DELETE")

	r.Handle("/v1/user", AuthMiddleware(http.HandlerFunc(showUserHandler))).Methods("GET")

	corsHandler := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		AllowCredentials: true,
	})

	http.Handle("/", handlers.LoggingHandler(os.Stdout, corsHandler.Handler(r)))

	log.Fatal(http.ListenAndServe(":8080", nil))
}
