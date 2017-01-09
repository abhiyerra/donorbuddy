package main

import (
	"context"
	"encoding/json"
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

	var fbUser struct {
		Name string `json:"name"`
		ID   string `json:"id"`
	}

	err = json.Unmarshal(response, &fbUser)
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	u := &User{FacebookID: fbUser.ID}
	_, err = config.DB.Model(u).Where("facebook_id = ?", fbUser.ID).OnConflict("DO NOTHING").SelectOrInsert()
	if err != nil {
		log.Println(err)
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	log.Println(*u)

	session, _ := config.Store.Get(r, "user")
	session.Values["ID"] = u.Id
	session.Save(r, w)

	w.Header().Set("Accept", "text/plain")
	http.Redirect(w, r, r.Referer(), http.StatusTemporaryRedirect)
}

func getSessionID(r *http.Request) int64 {
	session, _ := config.Store.Get(r, "user")
	val := session.Values["ID"]
	userID, ok := val.(int64)
	if !ok {
		return 0
	}
	return userID
}

const (
	UserKey = "User"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			userID = getSessionID(r)
			user   = User{Id: userID}
		)

		if userID == 0 {
			respondJson(w, r, struct{ Error string }{"Need to login"})
			return
		}

		log.Println(user)

		if err := config.DB.Select(&user); err != nil {
			respondJson(w, r, err)
			return
		}

		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), UserKey, user)))
	})
}

func UserValue(r *http.Request) User {
	return r.Context().Value(UserKey).(User)
}
