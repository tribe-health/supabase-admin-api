package optimizations

import "github.com/sirupsen/logrus"

type PgBouncerSettings struct {
	MaxClientConn   int `conf:"max_client_conn"`
	DefaultPoolSize int `conf:"default_pool_size"`
}

var (
	PgBouncerRecommendations = map[InstanceType]PgBouncerSettings{
		"t4g.micro": {
			MaxClientConn:   100,
			DefaultPoolSize: 5,
		},
		"t4g.small": {
			MaxClientConn:   200,
			DefaultPoolSize: 5,
		},
		"t4g.medium": {
			MaxClientConn:   400,
			DefaultPoolSize: 5,
		},
		"m6g.medium": {
			MaxClientConn:   400,
			DefaultPoolSize: 5,
		},
		"m6g.large": {
			MaxClientConn:   800,
			DefaultPoolSize: 5,
		},
		"m6g.xlarge": {
			MaxClientConn:   1600,
			DefaultPoolSize: 10,
		},
		"m6g.2xlarge": {
			MaxClientConn:   3200,
			DefaultPoolSize: 15,
		},
		"m6g.4xlarge": {
			MaxClientConn:   6400,
			DefaultPoolSize: 25,
		},
		"m6g.8xlarge": {
			MaxClientConn:   12800,
			DefaultPoolSize: 45,
		},
		"m6g.12xlarge": {
			MaxClientConn:   12800,
			DefaultPoolSize: 60,
		},
		"m6g.16xlarge": {
			MaxClientConn:   12800,
			DefaultPoolSize: 70,
		},
	}
)

func OptimizePgBouncer(destinationFile string, instanceType InstanceType) error {
	settings, ok := PgBouncerRecommendations[instanceType]
	if !ok {
		logrus.WithField("instanceType", instanceType).Warn("Using fallback recommendations.")
		settings, _ = PgBouncerRecommendations[FallbackInstanceType]
	}
	return writeRecommendationsToFile(settings, destinationFile)
}
