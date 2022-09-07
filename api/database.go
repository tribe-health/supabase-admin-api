package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type UpgradeDatabaseConfiguration struct {
	Password     string `json:"password"`
	MajorVersion int    `json:"major_version"`
}

func (a *API) PrepareDatabaseUpgrade(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/prepare_pg_upgrade.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "couldn't unmount data volume")
	}
	logrus.WithField("output", string(output)).Info("Unmounted data volume")

	return nil
}

func (a *API) InitiateDatabaseUpgrade(w http.ResponseWriter, r *http.Request) error {
	params := &UpgradeDatabaseConfiguration{}

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo /root/initiate_pg_upgrade.sh %d %s", params.MajorVersion, params.Password))
	output, err := cmd.Output()
	if err != nil {
		logrus.WithField("output", string(output)).Info("errored")
		return errors.Wrap(err, "couldn't run pg_upgrade")
	}
	logrus.WithField("output", string(output)).Info("Finished running pg_upgrade")

	return nil
}

func (a *API) CompleteDatabaseUpgrade(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/complete_pg_upgrade.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "couldn't complete pg_upgrade")
	}
	logrus.WithField("output", string(output)).Info("Completed database upgrade")

	return nil
}
