package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"time"
)

type PostgrestCollector struct {
	up     *prometheus.Desc
	client *http.Client
	url    string
}

func NewPostgrestCollector(postgrestUrl string) *PostgrestCollector {
	return &PostgrestCollector{
		up: prometheus.NewDesc("postgrest_up", "PostgREST status", nil, nil),
		client: &http.Client{
			Timeout: 500 * time.Millisecond,
		},
		url: postgrestUrl,
	}
}

func (c *PostgrestCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
}

func (c *PostgrestCollector) Collect(ch chan<- prometheus.Metric) {
	resp, err := c.client.Head(c.url)
	status := float64(0)
	if err == nil && resp != nil && resp.StatusCode == 200 {
		status = float64(1)
	}
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, status)
}
