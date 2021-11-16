package api

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"

	"github.com/go-chi/chi"
)

const postgrestService string = "postgrest.service"

const pgListenService string = "pg_listen.service"

const gotrueService string = "gotrue.service"

const adminAPIService string = "admin-api.service"

const realtimeService string = "realtime.service"

const kongService string = "kong.service"
const kongErrorService string = "kong.service"

const sysService string = "services.slice"

// GetLogContents is the method for returning the contents of a given log file
func (a *API) GetLogContents(w http.ResponseWriter, r *http.Request) error {
	application := chi.URLParam(r, "application")

	// fetchType is head, tail
	fetchType := chi.URLParam(r, "type")

	// number of lines if head or tail
	n := chi.URLParam(r, "n")

	// default is concatenate entire file
	reverseArg := "-r"
	arg0 := "-n"
	arg1 := "100"
	serviceName := sysService

	switch application {
	case "test":
		serviceName = "./README.md"
	case "gotrue":
		serviceName = gotrueService
	case "postgrest":
		serviceName = postgrestService
	case "pglisten":
		serviceName = pgListenService
	case "kong":
		serviceName = kongService
	case "kong-error":
		serviceName = kongErrorService
	case "realtime":
		serviceName = realtimeService
	case "admin":
		serviceName = adminAPIService
	case "syslog":
		serviceName = sysService
	}

	switch fetchType {
	case "head":
		reverseArg = "--no-pager" // no-op
		arg0 = "-n"
		arg1 = n
	case "tail":
		reverseArg = "-r"
		arg0 = "-n"
		arg1 = n
	}

	cmd := exec.Command("journalctl", "-u", serviceName, reverseArg, arg0, arg1, "--no-pager")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	return sendJSON(w, http.StatusOK, string(stdout))
}
