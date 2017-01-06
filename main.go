package main

import (
	"net/http"
	"time"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/objx"

	"github.com/gorilla/mux"
	"gopkg.in/pg.v5"
)

type config struct {
	StripeSecretKey string
}

// Org is a table of the organizations that we will be supporting for
// donations..
type Org struct {
	Id int64

	Name string
	EIN  string

	Address string
	City    string
	State   string
	Country string

	Category string

	Verified bool

	StripeCustomerID string `json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Users   []User   `pg:",many2many:user_orgs"`
	Ledgers []Ledger `pg:",many2many:ledgers"`
}

// User contains the information about the user who will be doing the
// donations.
type User struct {
	Id int64

	FacebookID string

	StripeCustomerID     string `json:"-"`
	StripeSubscriptionID string `json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Orgs    []Org    `pg:",many2many:user_orgs"`
	Ledgers []Ledger `pg:",many2many:ledgers"`
}

// UserOrg creates a map of the org that the user is donating to. A
// user can have more than one of the same user and org map. That is
// how we make donating more to a certain cause easier.
type UserOrg struct {
	UserId int64 `sql:",pk"`
	OrgId  int64 `sql:",pk"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// PaymentLedger contains the leger of all the transations that have
// happened in the system. So we have a complete list of transactions
// as they have happened.
type Ledger struct {
	UserId int64 `sql:",pk"`
	OrgId  int64 `sql:",pk"`
	Amount float64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := gomniauth.Provider("facebook")
	if err != nil {
		panic(err)
	}

	state := gomniauth.NewState("after", "success")

	// This code borrowed from goweb example and not fixed.
	// if you want to request additional scopes from the provider,
	// pass them as login?scope=scope1,scope2
	//options := objx.MSI("scope", ctx.QueryValue("scope"))

	authUrl, err := provider.GetBeginAuthURL(state, nil)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect
	http.Redirect(w, r, authUrl, http.StatusFound)

}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := gomniauth.Provider("facebook")
	if err != nil {
		panic(err)
	}

	omap, err := objx.FromURLQuery(r.URL.RawQuery)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	creds, err := provider.CompleteAuth(omap)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	/*
		// This code borrowed from goweb example and not fixed.
		// get the state
		state, err := gomniauth.StateFromParam(ctx.QueryValue("state"))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		// redirect to the 'after' URL
		afterUrl := state.GetStringOrDefault("after", "error?e=No after parameter was set in the state")
	*/

	// load the user
	user, userErr := provider.GetUser(creds)

	if userErr != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	data := fmt.Sprintf("%#v", user)
	io.WriteString(w, data)

	// redirect
	//return goweb.Respond.WithRedirect(ctx, afterUrl)

}

func showOrgssHandler(w http.ResponseWriter, r *http.Request) {

}

func searchOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		orgs []Org
	)

	err := db.Model(&orgs).Where().Limit(50).Select()
	if err != nil {
		renderJson(w, r, err)
	}

	renderJson(w, r, orgs)
}

func showPaymentsHandler(w http.ResponseWriter, r *http.Request) {

}

func updatePaymentsHandler(w http.ResponseWriter, r *http.Request) {

}

func deletePaymentsHandler(w http.ResponseWriter, r *http.Request) {

}

func main() {
	gomniauth.SetSecurityKey("yLiCQYG7CAflDavqGH461IO0MHp7TEbpg6TwHBWdJzNwYod1i5ZTbrIF5bEoO3oP") // NOTE: DO NOT COPY THIS - MAKE YOR OWN!
	gomniauth.WithProviders(
		facebook.New("537611606322077", "f9f4d77b3d3f4f5775369f5c9f88f65e", "http://localhost:8080/auth/facebook/callback"),
	)

	r := mux.NewRouter()
	r.HandleFunc("/auth/facebook", loginHandler)
	r.HandleFunc("/auth/facebook/callback", loginCallbackHandler)

	// r.HandleFunc("/v1/orgs", ArticlesHandler).Methods("GET")
	// r.HandleFunc("/v1/orgs/{orgId}", ArticlesHandler).Methods("GET")
	// r.HandleFunc("/v1/orgs/search", ArticlesHandler).Methods("GET")

	// r.HandleFunc("/v1/payments", ArticlesHandler).Methods("POST")
	// r.HandleFunc("/v1/payments", ArticlesHandler).Methods("UPDATE")
	// r.HandleFunc("/v1/payments", ArticlesHandler).Methods("DELETE")
	// r.HandleFunc("/v1/payments/stripe-callback", ArticlesHandler)

	// r.HandleFunc("/v1/user", ArticlesHandler)
	// r.HandleFunc("/v1/user/orgs", ArticlesHandler).Methods("GET")
	// r.HandleFunc("/v1/user/orgs", ArticlesHandler).Methods("DELETE")
	// r.HandleFunc("/v1/user/ledgers", ArticlesHandler).Methods("GET")

	http.Handle("/", r)
}
