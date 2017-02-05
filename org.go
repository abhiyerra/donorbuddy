package main

import (
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/mux"
)

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

func showOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars  = mux.Vars(r)
		orgId = vars["orgId"]
		org   Org
		err   error
	)

	if org.Id, err = strconv.ParseInt(orgId, 10, 64); err != nil {
		log.Println(err)
		respondJson(w, r, err)
		return
	}

	if err = config.DB.Select(&org); err != nil {
		log.Println(err)
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, org)
}

func searchOrgsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		orgs   []Org
		search = config.DB.Model(&orgs)

		city     = r.FormValue("city")
		state    = r.FormValue("state")
		category = r.FormValue("category")
	)

	if city != "" {
		search = search.Where("city = ?", city)
	}

	if state != "" {
		search = search.Where("state = ?", state)
	}

	if category != "" {
		search = search.Where("category = ?", category)
	}

	// select ein, name, state, city from orgs where to_tsvector('english', name) @@ plainto_tsquery('english', 'chimera');

	err := search.Order("name asc").Limit(50).Select()
	if err != nil {
		respondJson(w, r, err)
	}

	respondJson(w, r, orgs)
}
