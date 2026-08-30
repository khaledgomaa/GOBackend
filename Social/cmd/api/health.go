// This file equivalent to Controller or fast api in .NET
package main

import "net/http"

func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Healthy"))
	w.WriteHeader(http.StatusOK)
}
