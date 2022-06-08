package api

import (
	"encoding/json"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
	"strconv"
)

// FileContents holds the content of a config file
type RestoreConfiguration struct {
	BackupName         string `json:"backup_name"`
	RecoveryTimeTarget string `json:"recovery_time_target"`
}

type BackupConfiguration struct {
	ProjectId int `json:"project_id"`
	BackupId  int `json:"backup_id"`
}

func (a *API) BackupDatabase(w http.ResponseWriter, r *http.Request) error {
	params := &BackupConfiguration{}

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		logrus.WithError(err).Warn("failed to decode parameters")
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	cmd := exec.Command("sudo", "/root/commence_walg_backup.sh", strconv.Itoa(params.ProjectId), strconv.Itoa(params.BackupId))
	output, err := cmd.Output()
	if err != nil {
		errMessage := "failed to execute WAL-G backup"
		logrus.WithField("output", string(output)).Warn(errMessage)
		return errors.Wrap(err, errMessage)
	}
	logrus.WithField("output", string(output)).Info("WAL-G backup completed")
	return nil
}

func (a *API) RestoreDatabase(w http.ResponseWriter, r *http.Request) error {
	params := &RestoreConfiguration{}

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		logrus.WithError(err).Warn("failed to decode parameters")
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	cmd := exec.Command("sudo", "/root/commence_walg_restore.sh", params.BackupName, params.RecoveryTimeTarget)
	output, err := cmd.Output()
	if err != nil {
		errMessage := "failed to execute WAL-G restore"
		logrus.WithField("output", string(output)).Warn(errMessage)
		return errors.Wrap(err, errMessage)
	}
	logrus.WithField("output", string(output)).Info("WAL-G restore completed")
	return nil
}

func (a *API) CompleteRestorationWALG(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/complete_walg_restore.sh")
	output, err := cmd.Output()
	if err != nil {
		errMessage := "failed to complete WAL-G restoration"
		logrus.WithField("output", string(output)).Warn(errMessage)
		return errors.Wrap(err, errMessage)
	}
	logrus.WithField("output", string(output)).Info("WAL-G restoration complete")
	return nil
}

func (a *API) EnableWALG(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/enable_walg.sh")
	output, err := cmd.Output()
	if err != nil {
		errMessage := "failed to enable WAL-G"
		logrus.WithField("output", string(output)).Warn(errMessage)
		return errors.Wrap(err, errMessage)
	}
	logrus.WithField("output", string(output)).Info("WAL-G enabled")
	return nil
}

func (a *API) DisableWALG(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/disable_walg.sh")
	output, err := cmd.Output()
	if err != nil {
		errMessage := "failed to disable WAL-G"
		logrus.WithField("output", string(output)).Warn(errMessage)
		return errors.Wrap(err, errMessage)
	}
	logrus.WithField("output", string(output)).Info("WAL-G disabled")
	return nil
}
