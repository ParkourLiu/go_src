package health

import (
	"net/http"
)

var HealthHandler func(w http.ResponseWriter, r *http.Request)

func init() {
	HealthHandler = func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(200)
		w.Write([]byte("ok"))
	}
}
