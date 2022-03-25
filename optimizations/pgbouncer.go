package optimizations

import "github.com/sirupsen/logrus"

type PgBouncerSettings struct {
	MaxClientConn   int `conf:"max_client_conn"`
	DefaultPoolSize int `conf:"default_pool_size"`
}

var (
	PgBouncerRecommendations = map[InstanceType]PgBouncerSettings{
		"t4g.micro": {
			MaxClientConn:   200,
			DefaultPoolSize: 15,
		},
		"t4g.small": {
			MaxClientConn:   200,
			DefaultPoolSize: 15,
		},
		"t4g.medium": {
			MaxClientConn:   200,
			DefaultPoolSize: 15,
		},
		"m6g.medium": {
			MaxClientConn:   200,
			DefaultPoolSize: 15,
		},
		"m6g.large": {
			MaxClientConn:   300,
			DefaultPoolSize: 15,
		},
		"m6g.xlarge": {
			MaxClientConn:   700,
			DefaultPoolSize: 20,
		},
		"m6g.2xlarge": {
			MaxClientConn:   1500,
			DefaultPoolSize: 25,
		},
		"m6g.4xlarge": {
			MaxClientConn:   3000,
			DefaultPoolSize: 32,
		},
		"m6g.8xlarge": {
			MaxClientConn:   6000,
			DefaultPoolSize: 64,
		},
		"m6g.12xlarge": {
			MaxClientConn:   9000,
			DefaultPoolSize: 96,
		},
		"m6g.16xlarge": {
			MaxClientConn:   12000,
			DefaultPoolSize: 128,
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
