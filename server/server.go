package server

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/raj-ptl/go-status-check/models"
	"github.com/raj-ptl/go-status-check/status"
)

func ServeRequests() {
	fmt.Println("Serving now...")
	http.HandleFunc("/", welcomeHandler)
	http.HandleFunc("/websites", statusHandler)
	http.ListenAndServe("127.0.0.1:9090", nil)
}

func welcomeHandler(w http.ResponseWriter, r *http.Request) {

	// Undefined Endpoint
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	fmt.Fprintf(w, "Welcome to the server\n")
}

func statusHandler(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")

	if r.Method == "GET" {
		var hc status.HttpChecker
		hc.Check(context.TODO(), "google.com")
		hc.Check(context.TODO(), "twitter.com")
		hc.Check(context.TODO(), "localhost:9090/")
		hc.Check(context.TODO(), "localhost:9090/undefined")
		jsonWIP, _ := json.Marshal("/GET WIP")
		w.Write(jsonWIP)
	} else if r.Method == "POST" {

		var sr models.StatusRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		errUnmarshal := decoder.Decode(&sr)

		var unmarshalErr *json.UnmarshalTypeError

		if errUnmarshal != nil {

			if errors.As(errUnmarshal, &unmarshalErr) {
				errorResponse(w, "Bad Request. Wrong Type provided for field "+unmarshalErr.Field, http.StatusBadRequest)
			} else {
				errorResponse(w, "Bad Request "+errUnmarshal.Error(), http.StatusBadRequest)
			}

		} else {
			srJson, errMarshal := json.Marshal(sr)
			if errMarshal != nil {
				jsonErrMarshal, _ := json.Marshal(errMarshal)
				w.Write(jsonErrMarshal)
			} else {
				w.Write(srJson)
			}
		}

	} else {
		jsonInvalidMethod, _ := json.Marshal("POST/GET Method expected on this endpoint")
		w.Write(jsonInvalidMethod)
	}

}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}