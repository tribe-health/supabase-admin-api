package monitors

type Monitor interface {
	StartMonitoring()
	StopMonitoring()
}
