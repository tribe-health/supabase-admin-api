package api

import (
	"encoding/json"
	"github.com/supabase/supabase-admin-api/api/upgrades"
	"net/http"
)

func (a *API) RequestFileDownload(w http.ResponseWriter, r *http.Request) error {
	var req upgrades.DownloadRequest
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&req); err != nil {
		return sendJSON(w, http.StatusBadRequest, err.Error())
	}
	response, err := a.upgrades.DownloadFile(&req)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}
	return sendJSON(w, http.StatusOK, response)
}
