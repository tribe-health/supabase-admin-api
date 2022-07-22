package api

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

type VolumeType = string

const (
	RootVolumeType VolumeType = "root"
	DataVolumeType VolumeType = "data"
)

type ExpandFileSystemConfiguration struct {
	VolumeType string `json:"volume_type"`
}

func ExpandFilesystem(w http.ResponseWriter, r *http.Request) error {
	params := &ExpandFileSystemConfiguration{}

	jsonDecoder := json.NewDecoder(r.Body)
	if err := jsonDecoder.Decode(params); err != nil {
		return sendJSON(w, http.StatusInternalServerError, err.Error())
	}

	if params.VolumeType != DataVolumeType && params.VolumeType != RootVolumeType {
		return sendJSON(w, http.StatusBadRequest, "Invalid value provided for `volume_type`.")
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo /root/grow_fs.sh %s", params.VolumeType))

	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "couldn't grow partition")
	}
	logrus.WithField("output", string(output)).Info("Resized disk")
	return nil
}
