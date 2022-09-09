package upgrades

import (
	"log"
	"os"
	"path"
	"path/filepath"

	"github.com/pkg/errors"
)

type UpgradeRequest struct {
	ServiceName                   string `json:"service_name"`
	BinarySourceAbsolutePath      string `json:"binary_source_absolute_path"`
	BinaryDestinationAbsolutePath string `json:"binary_destination_absolute_path"`
}

type UpgradeResponse struct {
	Status     string `json:"status"`
	BackupPath string `json:"backup_path"`
}

func (u *Upgrades) ExecuteUpgrade(r *UpgradeRequest) (*UpgradeResponse, error) {
	if _, err := os.Stat(r.BinaryDestinationAbsolutePath); err == nil {
		flName := filepath.Base(r.BinaryDestinationAbsolutePath)
		backupPath := path.Join(u.Config.DestinationDir, Transient, flName)
		log.Printf("moving %s -> %s\n", flName, backupPath)
		err := os.Rename(r.BinaryDestinationAbsolutePath, backupPath)
		if err != nil {
			return nil, errors.Wrapf(err, "failed to move old binary to backup %s -> %s", r.BinaryDestinationAbsolutePath, backupPath)
		}
	}
	err := os.Rename(r.BinarySourceAbsolutePath, r.BinaryDestinationAbsolutePath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to move new binary %s -> %s", r.BinarySourceAbsolutePath, r.BinaryDestinationAbsolutePath)
	}

	// TODO (darora): restart service

	return &UpgradeResponse{
		Status:     "",
		BackupPath: "",
	}, nil
}
