package main

import (
	"net/http"
)

// HealthCheck godoc
//
//	@Summary		Health Check
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Accept			json
//	@Produce		json
//
//	@Success		200	{object}	map[string]string
//	@Failure		500	{object}	error
//	@Security		ApiKeyAuth
//	@Router			/health [get]
func (app *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     app.config.env,
		"version": app.config.version,
	}

	if err := app.jsonResponse(w, http.StatusOK, data); err != nil {
		app.internalServerError(w, r, err)
	}

}
