package optimizations

import (
	"bytes"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/pkg/errors"
	"os"
	"reflect"
)

type ServerSettings struct {
	CheckpointCompletionTarget    float32 `postgres:"checkpoint_completion_target"`
	DefaultStatisticsTarget       int     `postgres:"default_statistics_target"`
	EffectiveCacheSize            string  `postgres:"effective_cache_size"`
	EffectiveIoConcurrency        int     `postgres:"effective_io_concurrency"`
	MaintenanceWorkMem            string  `postgres:"maintenance_work_mem"`
	MaxConnections                int     `postgres:"max_connections"`
	MaxParallelMaintenanceWorkers int     `postgres:"max_parallel_maintenance_workers"`
	MaxParallelWorkers            int     `postgres:"max_parallel_workers"`
	MaxParallelWorkersPerGather   int     `postgres:"max_parallel_workers_per_gather"`
	MaxWalSize                    string  `postgres:"max_wal_size"`
	MaxWorkerProcesses            int     `postgres:"max_worker_processes"`
	MinWalSize                    string  `postgres:"min_wal_size"`
	RandomPageCost                float32 `postgres:"random_page_cost"`
	SharedBuffers                 string  `postgres:"shared_buffers"`
	WalBuffers                    string  `postgres:"wal_buffers"`
	WorkMem                       string  `postgres:"work_mem"`
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

func generateSettings(settings *ServerSettings) (*string, error) {
	var buffer bytes.Buffer
	val := reflect.ValueOf(settings)
	for _, field := range reflect.VisibleFields(reflect.TypeOf(*settings)) {
		serverFieldName := field.Tag.Get("postgres")
		f := val.Elem().FieldByName(field.Name)
		if f.IsZero() {
			continue
		}
		switch f.Kind() {
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			buffer.WriteString(fmt.Sprintf("%s = %d\n", serverFieldName, f.Int()))
			break
		case reflect.Float64, reflect.Float32:
			buffer.WriteString(fmt.Sprintf("%s = %0.03f\n", serverFieldName, f.Float()))
			break
		case reflect.String:
			buffer.WriteString(fmt.Sprintf("%s = '%s'\n", serverFieldName, f.String()))
			break
		default:
			return nil, fmt.Errorf("unsupported type encountered for field %+v: %+v", f, f.Kind())
		}
	}
	return aws.String(buffer.String()), nil
}

func writeRecommendationsToFile(settings *ServerSettings, destinationFilePath string) error {
	output, err := generateSettings(settings)
	if err != nil {
		return errors.Wrap(err, "couldn't serialize settings")
	}
	err = os.WriteFile(destinationFilePath, []byte(*output), 0644)
	if err != nil {
		return errors.Wrapf(err, "couldn't write recommendations to %s", destinationFilePath)
	}
	return nil
}

func OptimizeSettings(destinationFile string, instanceType InstanceType) error {
	settings, ok := ServerRecommendations[instanceType]
	if !ok {
		return fmt.Errorf("don't have recommended settings for instance type '%s'", instanceType)
	}
	return writeRecommendationsToFile(&settings, destinationFile)
}
