package main

import (
	"encoding/json"
	"fmt"
	"net/http"
)

func respondJson(w http.ResponseWriter, r *http.Request, i interface{}) {
	w.Header().Set("Content-Type", "application/json")

	jsonResponse, err := json.MarshalIndent(i, "", "\t")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	callback := r.FormValue("callback")

	callbackPrefix := ""
	if callback != "" {
		callbackPrefix = callback + "("
	}

	callbackSuffix := ""
	if callback != "" {
		callbackSuffix = ")"
	}

	switch v := i.(type) {
	case error:
		http.Error(w, callbackPrefix+"{\"Error\":"+string(v)+"}"+callbackSuffix, 404)
		return
	}

	fmt.Fprintf(w, callbackPrefix+string(jsonResponse)+callbackSuffix)
}
