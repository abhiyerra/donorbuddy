package main

import (
	"net/http"
	"time"

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

	StripeCustomerID string

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

	StripeCustomerID     string
	StripeSubscriptionID string

	CreatedAt time.Time
	UpdatedAt time.Time

	Orgs    []Org    `pg:",many2many:user_orgs"`
	Ledgers []Ledger `pg:",many2many:ledgers"`
}

// UserOrg creates a map of the org that the user is donating to. A
// user can have more than one of the same user and org map. That is
// how we make donating more to a certain cause easier.
type UserOrg struct {
	UserID int64 `sql:",pk"`
	OrgID  int64 `sql:",pk"`

	CreatedAt time.Time
	UpdatedAt time.Time
}

// PaymentLedger contains the leger of all the transations that have
// happened in the system. So we have a complete list of transactions
// as they have happened.
type Ledger struct {
	UserID int64 `sql:",pk"`
	OrgID  int64 `sql:",pk"`
	Amount float64

	CreatedAt time.Time
	UpdatedAt time.Time
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", HomeHandler)
	r.HandleFunc("/products", ProductsHandler)
	r.HandleFunc("/articles", ArticlesHandler)
	http.Handle("/", r)

}
