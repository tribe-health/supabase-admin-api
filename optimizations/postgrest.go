package optimizations

import (
	"github.com/sirupsen/logrus"
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
		logrus.WithField("instanceType", instanceType).Warn("Using fallback recommendations.")
		settings, _ = PostgrestServerRecommendations[FallbackInstanceType]
	}
	return writeRecommendationsToFile(settings, destinationFile)
}
