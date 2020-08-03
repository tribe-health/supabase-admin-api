package api

import (
	"net/http"
)

// GetKongYaml is the endpoint for fetching current goauth email config
func (a *API) GetKongYaml(w http.ResponseWriter, r *http.Request) error {

	return sendJSON(w, http.StatusOK, 0)
}

// SetKongYaml is the endpoint for fetching current goauth email config
func (a *API) SetKongYaml(w http.ResponseWriter, r *http.Request) error {

	return sendJSON(w, http.StatusOK, 0)
}
