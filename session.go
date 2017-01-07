package main

import (
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/stretchr/gomniauth"
	"github.com/stretchr/objx"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	provider, err := gomniauth.Provider("facebook")
	if err != nil {
		log.Println(err)
	}

	authURL, err := provider.GetBeginAuthURL(gomniauth.NewState("after", "success"), nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// redirect
	http.Redirect(w, r, authURL, http.StatusFound)
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
