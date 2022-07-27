package api

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/sirupsen/logrus"
	metrics "github.com/supabase/supabase-admin-api/api/metrics_endpoint"
)

func upstreamServers() (*httptest.Server, *httptest.Server) {
	server1 := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				_, _ = fmt.Fprint(w, `# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds{k="v"} 1024
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.0584064e+07
`)
			}))

	server2 := httptest.NewServer(
		http.HandlerFunc(
			func(w http.ResponseWriter, req *http.Request) {
				_, _ = fmt.Fprint(w, `# HELP node_memory_KernelStack_bytes Memory information field KernelStack_bytes.
# TYPE node_memory_KernelStack_bytes gauge
node_memory_KernelStack_bytes 3.887104e+06
# HELP node_memory_Mapped_bytes Memory information field Mapped_bytes.
# TYPE node_memory_Mapped_bytes gauge
node_memory_Mapped_bytes 2.45776384e+08
`)
			}))
	return server1, server2
}

func TestUpstreamMetricsEndpoint(t *testing.T) {
	s1, s2 := upstreamServers()
	defer s1.Close()
	defer s2.Close()
	client := http.Client{
		Timeout: 1 * time.Second,
	}
	var parser expfmt.TextParser
	metricsProvider := metrics.Metrics{
		Sources: []metrics.MetricsSource{{
			Config: metrics.MetricsSourceConfig{
				Name: "db_system_metrics",
				Url:  s1.URL,
				LabelsToAttach: []*io_prometheus_client.LabelPair{{
					Name:  aws.String("project"),
					Value: aws.String("12345"),
				}, {
					Name:  aws.String("Name"),
					Value: aws.String("prod-db-ref"),
				}},
			},
			HttpClient: &client,
			Logger:     logrus.New(),
			Parser:     &parser,
		}, {
			Config: metrics.MetricsSourceConfig{
				Name: "middleware_system_metrics",
				Url:  s2.URL,
				LabelsToAttach: []*io_prometheus_client.LabelPair{{
					Name:  aws.String("project"),
					Value: aws.String("12345"),
				}, {
					Name:  aws.String("Name"),
					Value: aws.String("prod-1-ref"),
				}},
			},
			HttpClient: &client,
			Logger:     logrus.New(),
			Parser:     &parser,
		}},
	}
	result := strings.Split(metricsProvider.GetMergedMetrics(), "\n")
	sort.Strings(result)
	expectedResult := strings.Split(`# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds{project="12345",Name="prod-db-ref",k="v"} 1024
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes{project="12345",Name="prod-db-ref"} 1.0584064e+07
# HELP node_memory_KernelStack_bytes Memory information field KernelStack_bytes.
# TYPE node_memory_KernelStack_bytes gauge
node_memory_KernelStack_bytes{project="12345",Name="prod-1-ref"} 3.887104e+06
# HELP node_memory_Mapped_bytes Memory information field Mapped_bytes.
# TYPE node_memory_Mapped_bytes gauge
node_memory_Mapped_bytes{project="12345",Name="prod-1-ref"} 2.45776384e+08
`, "\n")
	sort.Strings(expectedResult)

	if !reflect.DeepEqual(result, expectedResult) {
		t.Fatalf("expected '%s' to equal '%s'", result, expectedResult)
	}
}
