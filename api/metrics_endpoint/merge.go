package metrics_endpoint

import (
	"bytes"
	prom "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/sirupsen/logrus"
	"io"
	"net/http"
)

type MetricsSourceConfig struct {
	Name           string            `yaml:"name"`
	Url            string            `yaml:"url"`
	LabelsToAttach []*prom.LabelPair `yaml:"labels_to_attach"`
	SkipTlsVerify  bool              `yaml:"skip_tls_verify" required:"false"`
	SourceTimeout  string            `yaml:"source_timeout" required:"false"`
}

type MetricsSource struct {
	Config     MetricsSourceConfig
	HttpClient *http.Client
	Logger     logrus.FieldLogger
	Parser     *expfmt.TextParser
}

type Metrics struct {
	Sources []MetricsSource
}

func (m *Metrics) GetMergedMetrics() string {
	var buffer bytes.Buffer
	for _, source := range m.Sources {
		buffer.Write(source.GetAndLabelMetrics())
	}
	return buffer.String()
}

func (s *MetricsSource) GetAndLabelMetrics() []byte {
	req, err := http.NewRequest("GET", s.Config.Url, nil)
	if err != nil {
		s.Logger.WithError(err).Warn("failed to create request")
		return []byte{}
	}
	req.Header.Set("Accept", "text/plain")
	resp, err := s.HttpClient.Do(req)
	if err != nil {
		s.Logger.WithError(err).Info("Failed to fetch upstream source")
		return []byte{}
	}
	return s.ParseAndLabelMetrics(resp.Body)
}

func (s *MetricsSource) ParseAndLabelMetrics(in io.Reader) []byte {
	var buffer bytes.Buffer
	mf, err := s.Parser.TextToMetricFamilies(in)
	if err != nil {
		s.Logger.WithError(err).Info("Failed to read upstream or parse metrics")
		return buffer.Bytes()
	}
	for _, v := range mf {
		for _, metric := range v.Metric {
			metric.Label = append(s.Config.LabelsToAttach, metric.Label...)
		}
		_, err := expfmt.MetricFamilyToText(&buffer, v)
		if err != nil {
			s.Logger.WithError(err).Info("Failed to write out metric family")
			return buffer.Bytes()
		}
	}
	return buffer.Bytes()
}
