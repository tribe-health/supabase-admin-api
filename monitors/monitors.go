package monitors

type MonitoringConfig struct {
	DiskUsage DiskUsageMonitorConfig `yaml:"disk_usage"`
}

type MonitorSet struct {
	diskUsage *DiskUsageMonitor
}

func NewMonitorSet(config MonitoringConfig) (*MonitorSet, error) {
	diskUsageMonitor, err := NewDiskUsageMonitor(config.DiskUsage)
	if err != nil {
		return nil, err
	}

	return &MonitorSet{
		diskUsage: diskUsageMonitor,
	}, nil
}

func (m *MonitorSet) StartMonitoring() {
	go m.diskUsage.StartMonitoring()
}

func (m *MonitorSet) StopMonitoring() {
	m.diskUsage.StopMonitoring()
}
