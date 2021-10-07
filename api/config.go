package api

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

const postgrestConfPath string = "/etc/postgrest.conf"
const postgrestConfPathOld string = "/etc/old.postgrest.conf"

const pgListenConfPath string = "/etc/pg_listen.conf"
const pgListenConfPathOld string = "/etc/old.pg_listen.conf"

const gotrueEnvPath string = "/etc/gotrue.env"
const gotrueEnvPathOld string = "/etc/old.gotrue.env"

// to avoid a breaking change that needs to be synchronized across various repos, we're accepting a param for this
// this can be removed in favour of a pre-defined const after a short transitionary period
const realtimeEnvPathPattern string = "/etc/%s.env"
const realtimeEnvPathOldPattern string = "/etc/old.%s.env"

const kongYmlPath string = "/etc/kong/kong.yml"
const kongYmlPathOld string = "/etc/kong/old.kong.yml"

const adminapiEnvPath string = "/etc/adminapi.env"
const adminapiEnvPathOld string = "/etc/old.adminapi.env"

// FileContents holds the content of a config file
type FileContents struct {
	RawContents     string `json:"raw_contents"`
	RestartServices bool   `json:"restart_services"`
}

// GetFileContents is the method for returning the contents of a given file
func (a *API) GetFileContents(w http.ResponseWriter, r *http.Request) error {
	var configFilePath string
	application := chi.URLParam(r, "application")

	switch application {
	case "test":
		configFilePath = "./README.md"
	case "gotrue":
		configFilePath = gotrueEnvPath
	case "postgrest":
		configFilePath = postgrestConfPath
	case "pglisten":
		configFilePath = pgListenConfPath
	case "kong":
		configFilePath = kongYmlPath
	case "realtime":
		configFilePath = fmt.Sprintf(realtimeEnvPathPattern, a.config.RealtimeServiceName)
	case "adminapi":
		configFilePath = adminapiEnvPath
	}

	contents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}
	fileContents := &FileContents{}
	fileContents.RawContents = string(contents)
	return sendJSON(w, http.StatusOK, fileContents)
}

// SetFileContents sets the data in a given file
func (a *API) SetFileContents(w http.ResponseWriter, r *http.Request) error {
	var configFilePath string
	var configFilePathOld string

	application := chi.URLParam(r, "application")

	switch application {
	case "test":
		configFilePath = "./README.md"
		configFilePathOld = "./old.README.md"
	case "gotrue":
		configFilePath = gotrueEnvPath
		configFilePathOld = gotrueEnvPathOld
	case "postgrest":
		configFilePath = postgrestConfPath
		configFilePathOld = postgrestConfPathOld
	case "pglisten":
		configFilePath = pgListenConfPath
		configFilePathOld = pgListenConfPathOld
	case "kong":
		configFilePath = kongYmlPath
		configFilePathOld = kongYmlPathOld
	case "realtime":
		configFilePath = fmt.Sprintf(realtimeEnvPathPattern, a.config.RealtimeServiceName)
		configFilePathOld = fmt.Sprintf(realtimeEnvPathOldPattern, a.config.RealtimeServiceName)
	case "adminapi":
		configFilePath = adminapiEnvPath
		configFilePathOld = adminapiEnvPathOld
	}

	params := &FileContents{}

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	err := os.Rename(configFilePath, configFilePathOld)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	f, err := os.Create(configFilePath)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	defer f.Close()

	bytesWritten, err := f.WriteString(params.RawContents)
	_ = bytesWritten
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}
	f.Sync()

	err = os.Chmod(configFilePath, 0664)
	if err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	if params.RestartServices {
		return a.RestartServices(w, r)
	}

	return sendJSON(w, http.StatusOK, map[string]int{"bytes_written": bytesWritten})
}
