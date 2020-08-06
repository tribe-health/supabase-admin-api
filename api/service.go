package api

import (
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
)

// RestartServices is the endpoint for fetching current goauth email config
func (a *API) RestartServices(w http.ResponseWriter, r *http.Request) error {
	sudo := "sudo"
	app := "systemctl"
	arg0 := "daemon-reload"

	cmd := exec.Command(sudo, app, arg0)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	fmt.Fprintf(os.Stdout, string(stdout))

	// need to do command as goroutine because adminapi gets killed and can't respond
	go func() {
		sudo := "sudo"
		app := "systemctl"
		arg0 := "restart"
		arg1 := "services.slice"

		// give time for http response
		time.Sleep(2 * time.Second)

		cmd = exec.Command(sudo, app, arg0, arg1)
		stdout, err = cmd.Output()

		if err != nil {
			fmt.Fprintf(os.Stderr, err.Error())
		}

		fmt.Fprintf(os.Stdout, string(stdout))
	}()

	return sendJSON(w, http.StatusOK, 200)
}

// RebootMachine is the endpoint for fetching current goauth email config
func (a *API) RebootMachine(w http.ResponseWriter, r *http.Request) error {
	// app := "reboot"
	// exec.Command(app)

	return sendJSON(w, http.StatusInternalServerError, "endpoint not yet implemented")
}
