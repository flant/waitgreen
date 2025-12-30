package wg

import (
	"encoding/json"
	"log"
	"net/http"
	"time"
	"waitgreen/modules/config"

	"github.com/uzhinskiy/lib.go/helpers"
)

type WaitGreen struct {
	conf config.Config
}

type apiRequest struct {
	WaitGreen bool `json:"waitgreen"`
}

var wgEnabled bool
var wgSetTime = time.Now()

var resetTime = 3 * time.Hour

func Run(cnf config.Config) {
	wg := WaitGreen{}
	wg.conf = cnf

	wgEnabled = cnf.App.DefaultWG
	go cleanup()

	http.HandleFunc("/", wg.ApiHandler)
	http.ListenAndServe(cnf.App.Bind+":"+cnf.App.Port, nil)
}

func (wg *WaitGreen) ApiHandler(w http.ResponseWriter, r *http.Request) {
	var request apiRequest

	defer r.Body.Close()
	remoteIP := helpers.GetIP(r.RemoteAddr, r.Header.Get("X-Real-IP"), r.Header.Get("X-Forwarded-For"))

	w.Header().Add("Access-Control-Allow-Origin", "*")
	w.Header().Add("Access-Control-Allow-Methods", "POST,OPTIONS,GET")
	w.Header().Add("Access-Control-Allow-Credentials", "true")
	w.Header().Add("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization")
	w.Header().Add("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Server", "wg")

	if r.Method == "OPTIONS" {
		return
	}

	switch r.Method {
	case http.MethodPost:
		{

			dec := json.NewDecoder(r.Body)
			dec.DisallowUnknownFields()

			err := dec.Decode(&request)
			if err != nil {
				httpJsonError(w, err.Error(), http.StatusInternalServerError)
				log.Println(remoteIP, "\t", r.Method, "\t", r.URL.Path, "\t", http.StatusInternalServerError, "\t", err.Error(), "\t", r.UserAgent())
				return
			}

			wgEnabled = request.WaitGreen
			wgSetTime = time.Now()

			resp := map[string]interface{}{
				"status": http.StatusOK,
			}
			j, _ := json.Marshal(resp)
			log.Println(remoteIP, "\t", r.Method, "\t", r.URL.Path, "\t", "200", "\t", r.UserAgent())
			w.Write(j)

		}
	case http.MethodGet:
		{
			resp := map[string]interface{}{
				"waitgreen": wgEnabled,
			}
			j, _ := json.Marshal(resp)
			log.Println(remoteIP, "\t", r.Method, "\t", r.URL.Path, "\t", "200", "\t", r.UserAgent())
			w.Write(j)
		}
	}

}

func httpJsonError(w http.ResponseWriter, errorText string, errorCode int) {
	w.Header().Set("X-Server", "wg")
	w.WriteHeader(errorCode)
	resp := map[string]interface{}{
		"status": errorCode,
		"error":  errorText,
	}
	j, _ := json.Marshal(resp)
	w.Write(j)
}

func cleanup() {
	for {
		now := time.Now()
		if !wgEnabled {
			if diff := now.Sub(wgSetTime); diff > resetTime {
				log.Printf("Set wgEnabled to true at %s", diff)
				wgEnabled = true
				wgSetTime = time.Now()
			}
		}
		// do some job
		time.Sleep(30 * time.Second)
	}
}
