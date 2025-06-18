package main

import (
	"net/http"
)

// healthcheckHandler godoc
//
//	@Summary		Healthcheck
//	@Description	Healthcheck endpoint
//	@Tags			ops
//	@Produce		json
//	@Success		200	{object}	string	"ok"
//	@Router			/health [get]
func (api *application) healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	data := map[string]string{
		"status":  "ok",
		"env":     api.config.env,
		"version": version,
	}

	if err := api.jsonResponse(w, http.StatusOK, data); err != nil {
		api.internalServerError(w, r, err)
	}
}
