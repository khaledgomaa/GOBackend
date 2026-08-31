package main

import (
	"net/http"

	"go.uber.org/zap"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Error("internal server error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusInternalServerError, "There is something went wrong")
}

func (app *application) badRequestError(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("bad request error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("not found error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusNotFound, "not found")
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warn("unauthorized error", zap.String("method", r.Method), zap.String("path", r.URL.Path), zap.Error(err))
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
