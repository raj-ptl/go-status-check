package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/raj-ptl/go-status-check/constants"
	"github.com/raj-ptl/go-status-check/models"
	"github.com/raj-ptl/go-status-check/status"
)

var WebsiteMap = status.ExposeMap()

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
		//var hc status.HttpChecker

		if len(status.WebsiteMap) == 0 {
			jsonMapNotInitialized, _ := json.Marshal(constants.NO_WEBSITES_ADDED)
			w.Write(jsonMapNotInitialized)
		} else {

			jsonResponse, errJsonResponseMarshal := json.Marshal(WebsiteMap)

			if errJsonResponseMarshal != nil {
				w.Write([]byte(errJsonResponseMarshal.Error()))
			} else {
				w.Write(jsonResponse)
			}

		}

	} else if r.Method == "POST" {

		var sr models.StatusRequest

		decoder := json.NewDecoder(r.Body)
		decoder.DisallowUnknownFields()
		errUnmarshal := decoder.Decode(&sr)

		var unmarshalErr *json.UnmarshalTypeError

		if errUnmarshal != nil {

			if errors.As(errUnmarshal, &unmarshalErr) {
				errorResponse(w, constants.BAD_REQUEST_UNKNOWN_FIELD+unmarshalErr.Field, http.StatusBadRequest)
			} else {
				errorResponse(w, constants.BAD_REQUEST+errUnmarshal.Error(), http.StatusBadRequest)
			}

		} else {
			_, errMarshal := json.Marshal(sr)
			if errMarshal != nil {
				jsonErrMarshal, _ := json.Marshal(errMarshal)
				w.Write(jsonErrMarshal)
			} else {
				//w.Write(srJson)
				jsonWIP, _ := json.Marshal("WIP")
				w.Write(jsonWIP)
			}

			ch := make(chan int)

			for _, site := range sr.Websites {

				// update Single Site
				go status.UpdateSingleSite(site, ch)

			}

		}

	} else {
		jsonInvalidMethod, _ := json.Marshal(constants.UNEXPECTED_ENDPOINT)
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
