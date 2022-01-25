package metrics_endpoint

import (
	"bytes"
	"github.com/aws/aws-sdk-go/aws"
	io_prometheus_client "github.com/prometheus/client_model/go"
	"github.com/prometheus/common/expfmt"
	"github.com/sirupsen/logrus"
	"reflect"
	"sort"
	"strings"
	"testing"
)

func TestMetricsSource_ParseAndLabelMetrics(t *testing.T) {
	buffer := bytes.NewBufferString(`
# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds{old_label="old value"} 1024
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes 1.0584064e+07
`)
	var parser expfmt.TextParser
	source := MetricsSource{
		Parser: &parser,
		Config: MetricsSourceConfig{
			Url: "",
			LabelsToAttach: []*io_prometheus_client.LabelPair{
				{
					Name:  aws.String("project"),
					Value: aws.String("8783"),
				},
				{
					Name:  aws.String("Name"),
					Value: aws.String("prod-1-abcdef"),
				},
			},
		},
		Logger: logrus.New(),
	}
	metrics := source.ParseAndLabelMetrics(buffer)

	expected := strings.Split(`# HELP process_max_fds Maximum number of open file descriptors.
# TYPE process_max_fds gauge
process_max_fds{project="8783",Name="prod-1-abcdef",old_label="old value"} 1024
# HELP process_resident_memory_bytes Resident memory size in bytes.
# TYPE process_resident_memory_bytes gauge
process_resident_memory_bytes{project="8783",Name="prod-1-abcdef"} 1.0584064e+07
`, "\n")
	relabeled := strings.Split(string(metrics), "\n")
	sort.Strings(relabeled)
	sort.Strings(expected)
	if !reflect.DeepEqual(relabeled, expected) {
		t.Fatalf("Failed to relabel metrics; %s != %s", relabeled, expected)
	}
}
