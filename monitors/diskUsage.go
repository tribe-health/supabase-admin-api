package monitors

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/sys/unix"
)

type DiskUsageMonitorConfig struct {
	Enabled              bool   `yaml:"enabled"`
	IntervalDuration     string `yaml:"interval_duration"`
	ReadOnlyModeTreshold int    `yaml:"readonly_mode_treshold"`
}

const DefaultDiskUsageMonitoringIntervalDuration = "5s"
const DefaultDatabaseDiskUsageReadOnlyTreshold = 97

type DiskUsageMonitor struct {
	enabled                       bool
	interval                      time.Duration
	readOnlyModeDiskUsageTreshold float64
	doneChan                      chan (bool)
	readOnlyModeEnabled           bool
	dataDiskPath                  string
}

func NewDiskUsageMonitor(config DiskUsageMonitorConfig) (*DiskUsageMonitor, error) {
	var dataDiskPath string
	_, err := os.Stat("/data")
	if os.IsNotExist(err) {
		dataDiskPath = "/"
	} else {
		dataDiskPath = "/data"
	}

	if config.IntervalDuration == "" {
		config.IntervalDuration = DefaultDiskUsageMonitoringIntervalDuration
	}

	if config.ReadOnlyModeTreshold == 0 {
		config.ReadOnlyModeTreshold = DefaultDatabaseDiskUsageReadOnlyTreshold
	}

	monitorDuration, err := time.ParseDuration(config.IntervalDuration)
	if err != nil {
		return nil, err
	}

	return &DiskUsageMonitor{
		enabled:                       config.Enabled,
		readOnlyModeDiskUsageTreshold: float64(config.ReadOnlyModeTreshold),
		interval:                      monitorDuration,
		doneChan:                      make(chan bool, 1),
		readOnlyModeEnabled:           false,
		dataDiskPath:                  dataDiskPath,
	}, nil
}

func (d *DiskUsageMonitor) IsEnabled() bool {
	return d.enabled
}

func (d *DiskUsageMonitor) StartMonitoring() {
	if !d.IsEnabled() {
		return
	}

	logrus.WithField("monitor", "disk usage").Infof("Starting disk usage monitor for path %s.", d.dataDiskPath)
	t := time.NewTicker(d.interval)

	for {
		select {
		case <-d.doneChan:
			logrus.WithField("monitor", "disk usage").Info("Received stop signal. Stopping disk usage monitor.")
			return
		case <-t.C:
			err := d.monitor()
			if err != nil {
				logrus.WithField("monitor", "disk usage").WithError(err).Error("Failed monitoring disk usage.")
			}
		}
	}
}

func (d *DiskUsageMonitor) monitor() error {
	var fsStats unix.Statfs_t

	err := unix.Statfs(d.dataDiskPath, &fsStats)
	if err != nil {
		return errors.Wrap(err, "Failed to get filesystem stats.")
	}

	usedBlocksPct := (1 - float64(fsStats.Bavail)/float64(fsStats.Blocks)) * 100

	if usedBlocksPct >= d.readOnlyModeDiskUsageTreshold {
		if d.readOnlyModeEnabled {
			return nil
		}
		err = d.setReadOnlyModeEnabled(true)
		if err != nil {
			return err
		}
	} else {
		if !d.readOnlyModeEnabled {
			return nil
		}
		err = d.setReadOnlyModeEnabled(false)
		if err != nil {
			return err
		}
	}

	return nil
}

func (d *DiskUsageMonitor) readOnlyModeOverrideExists() (bool, error) {
	cmd := exec.Command("/bin/sh", "-c", "sudo /root/manage_readonly_mode.sh check_override")

	output, err := cmd.Output()
	if err != nil {
		return false, errors.Wrapf(err, "Couldn't check readonly mode override presence: %s", output)
	}

	if string(output) == "1" {
		return true, nil
	}

	return false, nil
}

func (d *DiskUsageMonitor) setReadOnlyModeEnabled(enableReadOnlyMode bool) error {
	var mode string
	if enableReadOnlyMode {
		mode = "on"
	} else {
		mode = "off"
	}

	overrideExists, err := d.readOnlyModeOverrideExists()
	if err != nil {
		return err
	}
	if overrideExists {
		logrus.WithField("monitor", "disk usage").Infof("Readonly mode override active. Bailing on setting mode to %s", mode)
		return nil
	}

	cmd := exec.Command("/bin/sh", "-c", fmt.Sprintf("sudo /root/manage_readonly_mode.sh set %s", mode))

	output, err := cmd.Output()
	if err != nil {
		return errors.Wrapf(err, "Couldn't set readonly mode (default_transaction_read_only) to %s: %s", mode, output)
	}

	logrus.Infof("Set readonly mode (default_transaction_read_only) to %s", mode)

	d.readOnlyModeEnabled = enableReadOnlyMode

	return nil
}

func (d *DiskUsageMonitor) StopMonitoring() {
	if !d.IsEnabled() {
		return
	}
	d.doneChan <- true
}
