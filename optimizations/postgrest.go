package optimizations

import (
	"fmt"
)

type PostgrestServerSettings struct {
	DbPool int `conf:"db-pool"`
}

var (
	PostgrestServerRecommendations = map[InstanceType]PostgrestServerSettings{
		"t4g.micro": {
			DbPool: 20,
		},
		"t4g.small": {
			DbPool: 30,
		},
		"m6g.medium": {
			DbPool: 40,
		},
		"m6g.large": {
			DbPool: 50,
		},
		"m6g.xlarge": {
			DbPool: 60,
		},
		"m6g.2xlarge": {
			DbPool: 70,
		},
	}
)

func OptimizePostgrest(destinationFile string, instanceType InstanceType) error {
	settings, ok := PostgrestServerRecommendations[instanceType]
	if !ok {
		return fmt.Errorf("don't have recommended settings for instance type '%s'", instanceType)
	}
	return writeRecommendationsToFile(settings, destinationFile)
}
