package optimizations

import "testing"

func TestGenerateConfig(t *testing.T) {
	settings := ServerSettings{
		MaxConnections:             200,
		WalBuffers:                 "200MB",
		CheckpointCompletionTarget: 0.9,
		MaxWorkerProcesses:         0}
	result, _ := generateSettings(settings)
	if *result != `checkpoint_completion_target = 0.900
max_connections = 200
wal_buffers = '200MB'
` {
		t.Fatal(*result)
	}
}
