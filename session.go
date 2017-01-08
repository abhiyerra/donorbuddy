package main

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"

	"golang.org/x/oauth2"
)

func loginHandler(w http.ResponseWriter, r *http.Request) {
	u, err := url.Parse(config.OAuthConf.Endpoint.AuthURL)
	if err != nil {
		log.Fatal("Parse: ", err)
	}

	params := url.Values{}
	params.Add("client_id", config.OAuthConf.ClientID)
	params.Add("scope", strings.Join(config.OAuthConf.Scopes, " "))
	params.Add("redirect_uri", config.OAuthConf.RedirectURL)
	params.Add("response_type", "code")
	params.Add("state", config.Auth.SecurityKey)

	u.RawQuery = params.Encode()
	http.Redirect(w, r, u.String(), http.StatusTemporaryRedirect)
}

func loginCallbackHandler(w http.ResponseWriter, r *http.Request) {
	state := r.FormValue("state")
	if state != config.Auth.SecurityKey {
		log.Printf("invalid oauth state, expected '%s', got '%s'\n", config.Auth.SecurityKey, state)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	code := r.FormValue("code")

	token, err := config.OAuthConf.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Printf("config.OAuthConf.Exchange() failed with '%s'\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp, err := http.Get("https://graph.facebook.com/me?access_token=" + url.QueryEscape(token.AccessToken))
	if err != nil {
		log.Printf("Get: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	defer resp.Body.Close()

	response, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Printf("ReadAll: %s\n", err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Printf("parseResponseBody: %s\n", string(response))

	http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
}
