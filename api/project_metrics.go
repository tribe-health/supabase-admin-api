package api

import (
	"github.com/sirupsen/logrus"
	"net/http"
)

const PlaceholderCacheKey = "placeholder"

func (a *API) ServeUpstreamMetrics(metricsProvider func(interface{}) (interface{}, error)) func(w http.ResponseWriter, r *http.Request) error {
	return func(w http.ResponseWriter, r *http.Request) error {
		metrics, err := metricsProvider(PlaceholderCacheKey)
		if err != nil {
			logrus.WithError(err).Warn("failed to get upstream metrics")
			return err
		}
		w.Header().Set("Content-Type", "text/plain")
		_, err = w.Write([]byte(metrics.(string)))
		return err
	}
}
