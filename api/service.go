package api

import (
	"net/http"
)

// RestartServices is the endpoint for fetching current goauth email config
func (a *API) RestartServices(w http.ResponseWriter, r *http.Request) error {

	return sendJSON(w, http.StatusOK, 0)
}

// RebootMachine is the endpoint for fetching current goauth email config
func (a *API) RebootMachine(w http.ResponseWriter, r *http.Request) error {

	return sendJSON(w, http.StatusOK, 0)
}
