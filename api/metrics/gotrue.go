package metrics

import (
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type GotrueCollector struct {
	up     *prometheus.Desc
	client *http.Client
	url    string
}

func NewGotrueCollector(gotrueUrl string) *GotrueCollector {
	return &GotrueCollector{
		up: prometheus.NewDesc("gotrue_up", "GoTrue status", nil, nil),
		client: &http.Client{
			Timeout: 500 * time.Millisecond,
		},
		url: gotrueUrl,
	}
}

func (c *GotrueCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
}

func (c *GotrueCollector) Collect(ch chan<- prometheus.Metric) {
	resp, err := c.client.Get(c.url)
	status := float64(0)
	if err == nil && resp != nil && resp.StatusCode == 200 {
		status = float64(1)
	}
	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, status)
}
