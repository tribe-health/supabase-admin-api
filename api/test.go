package api

import (
	"net/http"
)

// TestGet returns a user
func (a *API) TestGet(w http.ResponseWriter, r *http.Request) error {
	ctx := r.Context()

	return sendJSON(w, http.StatusOK, "hello")
}
