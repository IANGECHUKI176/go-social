package main

import (
	"net/http"
)

func (app *application) internalServerError(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("internal server error: %s path:%s error %s", r.Method, r.URL.Path, err)
	app.logger.Errorw("internal server error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusInternalServerError, "The server encountered a problem")
}

func (app *application) badRequestResponse(w http.ResponseWriter, r *http.Request, err error) {
	app.logger.Warnf("bad request error: %s path:%s error %s", r.Method, r.URL.Path, err.Error())
	writeJSONError(w, http.StatusBadRequest, err.Error())
}

func (app *application) notFoundResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("not found error: %s path:%s error %s", r.Method, r.URL.Path, err)
	app.logger.Warnf("not found error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusNotFound, "The requested resource could not be found")
}

func (app *application) conflictResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("conflict error: %s path:%s error %s", r.Method, r.URL.Path, err)
	app.logger.Errorf("conflict error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusConflict, err.Error())
}

func (app *application) unauthorizedErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("unauthorized error: %s path:%s error %s", r.Method, r.URL.Path, err)
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}

func (app *application) unauthorizedBasicErrorResponse(w http.ResponseWriter, r *http.Request, err error) {
	//log.Printf("unauthorized error: %s path:%s error %s", r.Method, r.URL.Path, err)
	app.logger.Warnf("unauthorized error", "method", r.Method, "path", r.URL.Path, "error", err.Error())
	w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
	writeJSONError(w, http.StatusUnauthorized, "unauthorized")
}
func (app *application) forbiddenResponse(w http.ResponseWriter, r *http.Request) {
	//log.Printf("forbidden error: %s path:%s error %s", r.Method, r.URL.Path, err)
	app.logger.Warnw("forbidden error", "method", r.Method, "path", r.URL.Path)
	writeJSONError(w, http.StatusForbidden, "forbidden")
}

func (app *application) rateLimitExceedResponse(w http.ResponseWriter, r *http.Request, retryAfter string) {
	app.logger.Warnw("rate limit exceeded", "method", r.Method, "path", r.URL.Path, "error", retryAfter)
	w.Header().Set("Retry-After", retryAfter)
	writeJSONError(w, http.StatusTooManyRequests, "rate limit exceeded, retry after: "+retryAfter)
}
