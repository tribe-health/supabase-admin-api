package api

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"net/http"
	"os/exec"
)

func ExpandFilesystem(w http.ResponseWriter, r *http.Request) error {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/grow_fs.sh")
	output, err := cmd.Output()
	if err != nil {
		return errors.Wrap(err, "couldn't grow partition")
	}
	logrus.WithField("output", string(output)).Info("Resized disk")
	return nil
}
