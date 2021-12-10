package optimizations

import "github.com/sirupsen/logrus"

type ServerSettings struct {
	CheckpointCompletionTarget    float32 `conf:"checkpoint_completion_target"`
	DefaultStatisticsTarget       int     `conf:"default_statistics_target"`
	EffectiveCacheSize            string  `conf:"effective_cache_size"`
	EffectiveIoConcurrency        int     `conf:"effective_io_concurrency"`
	MaintenanceWorkMem            string  `conf:"maintenance_work_mem"`
	MaxConnections                int     `conf:"max_connections"`
	MaxParallelMaintenanceWorkers int     `conf:"max_parallel_maintenance_workers"`
	MaxParallelWorkers            int     `conf:"max_parallel_workers"`
	MaxParallelWorkersPerGather   int     `conf:"max_parallel_workers_per_gather"`
	MaxWalSize                    string  `conf:"max_wal_size"`
	MaxWorkerProcesses            int     `conf:"max_worker_processes"`
	MinWalSize                    string  `conf:"min_wal_size"`
	RandomPageCost                float32 `conf:"random_page_cost"`
	SharedBuffers                 string  `conf:"shared_buffers"`
	WalBuffers                    string  `conf:"wal_buffers"`
	WorkMem                       string  `conf:"work_mem"`
}

var (
	ServerRecommendations = map[InstanceType]ServerSettings{
		"t4g.micro": {
			CheckpointCompletionTarget:    0.9,
			DefaultStatisticsTarget:       100,
			EffectiveCacheSize:            "768MB",
			EffectiveIoConcurrency:        200,
			MaintenanceWorkMem:            "64MB",
			MaxConnections:                200,
			MaxParallelMaintenanceWorkers: 1,
			MaxParallelWorkers:            2,
			MaxParallelWorkersPerGather:   1,
			MaxWalSize:                    "4GB",
			MaxWorkerProcesses:            2,
			MinWalSize:                    "1GB",
			RandomPageCost:                1.1,
			SharedBuffers:                 "256MB",
			WalBuffers:                    "7864kB",
			WorkMem:                       "1310kB",
		},
		"t4g.small": {
			CheckpointCompletionTarget:    0.9,
			DefaultStatisticsTarget:       100,
			EffectiveCacheSize:            "1536MB",
			EffectiveIoConcurrency:        200,
			MaintenanceWorkMem:            "128MB",
			MaxConnections:                200,
			MaxParallelMaintenanceWorkers: 1,
			MaxParallelWorkers:            2,
			MaxParallelWorkersPerGather:   1,
			MaxWalSize:                    "4GB",
			MaxWorkerProcesses:            2,
			MinWalSize:                    "1GB",
			RandomPageCost:                1.1,
			SharedBuffers:                 "512MB",
			WalBuffers:                    "16MB",
			WorkMem:                       "2621kB",
		},
		"m6g.medium": {
			CheckpointCompletionTarget:    0.9,
			DefaultStatisticsTarget:       100,
			EffectiveCacheSize:            "3GB",
			EffectiveIoConcurrency:        200,
			MaintenanceWorkMem:            "256MB",
			MaxConnections:                200,
			MaxParallelMaintenanceWorkers: 1,
			MaxParallelWorkers:            2,
			MaxParallelWorkersPerGather:   1,
			MaxWalSize:                    "4GB",
			MaxWorkerProcesses:            2,
			MinWalSize:                    "1GB",
			RandomPageCost:                1.1,
			SharedBuffers:                 "1GB",
			WalBuffers:                    "16MB",
			WorkMem:                       "2621kB",
		},
		"m6g.large": {
			CheckpointCompletionTarget:    0.9,
			DefaultStatisticsTarget:       100,
			EffectiveCacheSize:            "6GB",
			EffectiveIoConcurrency:        200,
			MaintenanceWorkMem:            "512MB",
			MaxConnections:                200,
			MaxParallelMaintenanceWorkers: 1,
			MaxParallelWorkers:            2,
			MaxParallelWorkersPerGather:   1,
			MaxWalSize:                    "4GB",
			MaxWorkerProcesses:            2,
			MinWalSize:                    "1GB",
			RandomPageCost:                1.1,
			SharedBuffers:                 "2GB",
			WalBuffers:                    "16MB",
			WorkMem:                       "10485kB",
		},
		"m6g.xlarge": {
			CheckpointCompletionTarget:    0.9,
			DefaultStatisticsTarget:       100,
			EffectiveCacheSize:            "12GB",
			EffectiveIoConcurrency:        200,
			MaintenanceWorkMem:            "1GB",
			MaxConnections:                200,
			MaxParallelMaintenanceWorkers: 2,
			MaxParallelWorkers:            4,
			MaxParallelWorkersPerGather:   2,
			MaxWalSize:                    "4GB",
			MaxWorkerProcesses:            4,
			MinWalSize:                    "1GB",
			RandomPageCost:                1.1,
			SharedBuffers:                 "4GB",
			WalBuffers:                    "16MB",
			WorkMem:                       "10485kB",
		},
		"m6g.2xlarge": {
			CheckpointCompletionTarget:    0.9,
			DefaultStatisticsTarget:       100,
			EffectiveCacheSize:            "24GB",
			EffectiveIoConcurrency:        200,
			MaintenanceWorkMem:            "2GB",
			MaxConnections:                200,
			MaxParallelMaintenanceWorkers: 4,
			MaxParallelWorkers:            8,
			MaxParallelWorkersPerGather:   4,
			MaxWalSize:                    "4GB",
			MaxWorkerProcesses:            8,
			MinWalSize:                    "1GB",
			RandomPageCost:                1.1,
			SharedBuffers:                 "8GB",
			WalBuffers:                    "16MB",
			WorkMem:                       "10485kB",
		},
	}
)

func OptimizePostgres(destinationFile string, instanceType InstanceType) error {
	settings, ok := ServerRecommendations[instanceType]
	if !ok {
		logrus.WithField("instanceType", instanceType).Warn("Using fallback recommendations.")
		settings, _ = ServerRecommendations[FallbackInstanceType]
	}
	return writeRecommendationsToFile(&settings, destinationFile)
}
