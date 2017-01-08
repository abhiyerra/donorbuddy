package main

import (
	"log"
	"net/http"
	"time"
)

// User contains the information about the user who will be doing the
// donations.
type User struct {
	Id int64

	FacebookID string

	StripeCustomerID     string `json:"-"`
	StripeSubscriptionID string `json:"-"`

	CreatedAt time.Time
	UpdatedAt time.Time

	Orgs    []*Org    `pg:",many2many:user_orgs"`
	Ledgers []*Ledger `pg:",many2many:ledgers"`
}

func showUserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		user = UserValue(r)
	)

	err := config.DB.Model(&user).Column("user.*", "Orgs").Select()
	//.Column("ledgers.*", "Ledgers")
	if err != nil {
		log.Println(err)
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, user)
}
