package api

import (
	"encoding/json"
	"net/http"
)

type RetrieveBans struct {
	Jails []string `json:"jails"`
}

type DeleteBans struct {
	Jails     []string `json:"jails"`
	IpAddress string   `json:"ip_address"`
}

func (a *API) GetCurrentBans(w http.ResponseWriter, r *http.Request) error {
	var req RetrieveBans
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&req); err != nil {
		return sendJSON(w, http.StatusBadRequest, err.Error())
	}

	ips, err := a.networkBans.ListBannedIps(req.Jails)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	return sendJSON(w, http.StatusOK, ips)
}

func (a *API) UnbanIp(w http.ResponseWriter, r *http.Request) error {
	var req DeleteBans
	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(&req); err != nil {
		return sendJSON(w, http.StatusBadRequest, err.Error())
	}

	err := a.networkBans.UnbanIp(req.Jails, req.IpAddress)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}
	return sendJSON(w, http.StatusOK, "ok")
}
