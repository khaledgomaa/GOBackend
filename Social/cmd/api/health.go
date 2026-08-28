package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy"))
	w.WriteHeader(http.StatusOK)
}
