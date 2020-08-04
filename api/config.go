package api

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/go-chi/chi"
)

const postgrestConfPath string = "./README.md"
const postgrestConfPathOld string = "/etc/old.postgrest.conf"

const pgListenConfPath string = "/etc/pg_listen.conf"
const pgListenConfPathOld string = "/etc/old.pg_listen.conf"

const gotrueEnvPath string = "/etc/gotrue.env"
const gotrueEnvPathOld string = "/etc/old.gotrue.env"

const realtimeEnvPath string = "/etc/supabase.env"
const realtimeEnvPathOld string = "/etc/old.supabase.env"

const kongYmlPath string = "/etc/kong/kong.yml"
const kongYmlPathOld string = "/etc/kong/old.kong.yml"

var configFilePath string
var configFilePathOld string

// FileContents holds the content of a config file
type FileContents struct {
	RawContents string `json:"raw_contents"`
}

// GetFileContents is the method for returning the contents of a given file
func (a *API) GetFileContents(w http.ResponseWriter, r *http.Request) error {
	application := chi.URLParam(r, "application")

	switch application {
	case "test":
		configFilePath = "./README.md"
	case "goauth":
		configFilePath = gotrueEnvPath
	case "postgrest":
		configFilePath = postgrestConfPath
	case "pglisten":
		configFilePath = pgListenConfPath
	case "kong":
		configFilePath = kongYmlPath
	case "realtime":
		configFilePath = realtimeEnvPath
	}

	contents, err := ioutil.ReadFile(configFilePath)
	if err != nil {
		return unprocessableEntityError(configFilePath)
	}
	fileContents := &FileContents{}
	fileContents.RawContents = string(contents)
	return sendJSON(w, http.StatusOK, fileContents)
}

// SetFileContents sets the data in a given file
func (a *API) SetFileContents(w http.ResponseWriter, r *http.Request) error {
	application := chi.URLParam(r, "application")

	switch application {
	case "test":
		configFilePath = "./README.md"
		configFilePathOld = "./old.README.md"
	case "goauth":
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
		configFilePath = realtimeEnvPath
		configFilePathOld = realtimeEnvPathOld
	}

	params := &FileContents{}

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		return badRequestError("Could not read config file data: %v", err)
	}

	err := os.Rename(configFilePath, configFilePathOld)
	if err != nil {
		return err
	}

	f, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	bytesWritten, err := f.WriteString(params.RawContents)
	if err != nil {
		return err
	}
	f.Sync()

	return sendJSON(w, http.StatusOK, map[string]int{"bytes_written": bytesWritten})
}
