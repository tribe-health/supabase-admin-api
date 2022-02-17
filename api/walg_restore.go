package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
)

// FileContents holds the content of a config file
type RestoreConfiguration struct {
	BackupName     string `json:"backup_name"`
	RecoveryTargetTime string   `json:"recovery_target_time"`
}

func RestoreDatabase(w http.ResponseWriter, r *http.Request) error {
	params := &RestoreConfiguration{}

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	var cmd *exec.Cmd
	if len(params.RecoveryTargetTime) == 0 {
		cmd = exec.Command("/bin/sh", "-c", "sudo /root/commence_walg_restore.sh", params.BackupName, params.RecoveryTargetTime)
	} else{
		cmd = exec.Command("/bin/sh", "-c", "sudo /root/commence_walg_restore.sh", params.BackupName)
	}
	output, err := cmd.Output()

	if err != nil {
		return errors.Wrap(err, "failed to execute WAL-G restoration")
	}
	logrus.WithField("output", string(output)).Info("WAL-G restoration completed")
	return nil
}
