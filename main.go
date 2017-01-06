package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/gomniauth/providers/facebook"
	"github.com/stretchr/objx"

	"github.com/stripe/stripe-go"
	"github.com/stripe/stripe-go/customer"
	"github.com/stripe/stripe-go/sub"

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

func showOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		orgId = vars["orgId"]
		org   Org
		err   error
	)

	if org.Id, err = strconv.ParseInt(orgId, 10, 64); err != nil {
		respondJson(w, r, err)
		return
	}

	if err = config.DB.Select(&org); err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, org)
}

func searchOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		orgs []Org
	)

	err := config.DB.Model(&orgs).Where().Limit(50).Select()
	if err != nil {
		respondJson(w, r, err)
	}

	respondJson(w, r, orgs)
}

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

	renderJson(w, r, struct{}{})
}

func deletePaymentsHandler(w http.ResponseWriter, r *http.Request) {
	err := sub.Cancel(
		"sub_9sed4J2K4jurwS", // TODO: Cancel user's subscription
	)

	renderJson(w, r, struct{}{})
}

func putUserOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars    = mux.Vars(r)
		orgId   = vars["orgId"]
		userOrg UserOrg
		err     error
	)

	if userOrg.OrgId, err = strconv.ParseInt(orgId, 10, 64); err != nil {
		respondJson(w, r, err)
		return
	}

	if err = config.DB.Insert(&userOrg); err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, userOrg)
}

func deleteUserOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars    = mux.Vars(r)
		orgId   = vars["orgId"]
		userOrg UserOrg
		err     error
	)

	if userOrg.OrgId, err = strconv.ParseInt(orgId, 10, 64); err != nil {
		respondJson(w, r, err)
		return
	}

	_, err = config.DB.Model(&userOrg).Where("user_id = ?user_id and org_id = ?org_id").Limit(1).Delete()
	if err != nil {
		respondJson(w, r, err)
		return

	}

	respondJson(w, r, struct{}{})
}

func showUserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars   = mux.Vars(r)
		userId = vars["userId"]
		user   User
		err    error
	)

	if user.Id, err = strconv.ParseInt(userId, 10, 64); err != nil {
		respondJson(w, r, err)
		return
	}

	if err = config.DB.Select(&user).Column("orgs.*", "Orgs").Column("ledgers.*", "Ledgers"); err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, user)
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
