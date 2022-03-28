package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"fmt"
)

// FileContents holds the content of a config file
type BackupConfiguration struct {
	ProjectId int `json:"project_id"`
	BackupId int   `json:"backup_id"`
}

func BackupDatabase(w http.ResponseWriter, r *http.Request) error {
	params := &BackupConfiguration{}

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		logrus.WithError(err).Warn("failed to decode parameters")
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

    completedCommand := fmt.Sprintf("%s %d %d", "sudo /root/commence_walg_backup.sh", params.ProjectId, params.BackupId)

	cmd := exec.Command("/bin/sh", "-c", completedCommand)
	output, err := cmd.Output()
	if err != nil {
		errMessage := "failed to execute WAL-G backup"
		logrus.WithError(output).Warn(errMessage)
		return errors.Wrap(err, errMessage)
	}
	logrus.WithField("output", string(output)).Info("WAL-G backup completed")
	return nil
}
