package main

import (
	"net/http"
	"strconv"
	"time"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
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

	Orgs    []Org    `pg:",many2many:user_orgs"`
	Ledgers []Ledger `pg:",many2many:ledgers"`
}

func showUserHandler(w http.ResponseWriter, r *http.Request) {
	var (
		vars   = mux.Vars(r)
		userID = vars["userId"]
		user   User
		err    error
	)

	if user.Id, err = strconv.ParseInt(userID, 10, 64); err != nil {
		respondJson(w, r, err)
		return
	}

	if err = config.DB.Model(&user).Column("orgs.*", "Orgs").Column("ledgers.*", "Ledgers").Select(); err != nil {
		respondJson(w, r, err)
		return
	}

	respondJson(w, r, user)
}
