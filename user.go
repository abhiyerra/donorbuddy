package main

import (
	"net/http"
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
