package api

import (
	"net/http"
)

// GetPGListenEnv is the endpoint for fetching current goauth email config
func (a *API) GetPGListenEnv(w http.ResponseWriter, r *http.Request) error {

	return sendJSON(w, http.StatusOK, 0)
}

// SetPGListenEnv is the endpoint for fetching current goauth email config
func (a *API) SetPGListenEnv(w http.ResponseWriter, r *http.Request) error {

	return sendJSON(w, http.StatusOK, 0)
}
