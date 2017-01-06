package main

import (
	"net/http"
)

func respondJson(w http.ResponseWriter, r *http.Request, i interface{}) {
	w.Header().Set("Content-Type", "application/json")

	switch i.(type) {
	case error:
		i = struct {
			Error string
		}{
			i.(error).String(),
		}
	}

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
	case ErrorResponse:
		http.Error(w, callbackPrefix+string(jsonResponse)+callbackSuffix, v.Code)
		return
	}

	fmt.Fprintf(w, callbackPrefix+string(jsonResponse)+callbackSuffix)
}
