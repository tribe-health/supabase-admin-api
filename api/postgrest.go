package api

import (
	"net/http"
)

// GetPostgrestEnv is the endpoint for fetching current goauth email config
func (a *API) GetPostgrestEnv(w http.ResponseWriter, r *http.Request) error {

	return sendJSON(w, http.StatusOK, 0)
}

// SetPostgrestEnv is the endpoint for fetching current goauth email config
func (a *API) SetPostgrestEnv(w http.ResponseWriter, r *http.Request) error {

	return sendJSON(w, http.StatusOK, 0)
}
