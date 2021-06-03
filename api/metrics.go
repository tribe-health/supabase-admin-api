package api

import (
	"bufio"
	"context"
	"fmt"
	"github.com/coreos/go-systemd/dbus"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var (
	reParens = regexp.MustCompile(`\((.*)\)`)
	namespace = "supabase"
)

type Metrics struct {
	registry      *prometheus.Registry
	meminfoFields map[string]prometheus.Gauge
	rtimeMetrics  map[string]func(interface{})
}

func NewMetrics() (*Metrics, error) {
	registry := prometheus.NewRegistry()
	memTotal := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name:      "memory_total_bytes",
	})
	memAvailable := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name: "memory_available_bytes",
	})
	meminfo := map[string]prometheus.Gauge{
		"MemTotal": memTotal,
		"MemAvailable": memAvailable,
	}
	for _, gauge := range meminfo {
		err := registry.Register(gauge)
		if err != nil {
			return nil, err
		}
	}

	rtimeRestarts := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name: "realtime_restarts_total",
	})
	rtimeMemory := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Name: "realtime_memory_bytes",
	})
	rtimeMetrics := map[string]func(interface{}){
		"NRestarts": func(val interface{}) {
			rtimeRestarts.Set(float64(val.(uint32)))
		},
		"MemoryCurrent": func(val interface{})  {
			rtimeMemory.Set(float64(val.(uint64)))
		},
	}
	for _, gauge := range []prometheus.Gauge{rtimeRestarts, rtimeMemory} {
		err := registry.Register(gauge)
		if err != nil {
			return nil, err
		}
	}
	return &Metrics{registry: registry, meminfoFields: meminfo, rtimeMetrics: rtimeMetrics}, nil
}

func (m *Metrics) UpdateRealtimeMetrics() error {
	ctx := context.Background()
	conn, err := dbus.NewSystemConnectionContext(ctx); if err != nil {
		return err
	}
	defer conn.Close()
	for key, consumer := range m.rtimeMetrics {
		val, err := conn.GetServicePropertyContext(ctx, "supabase.service", key); if err != nil {
			return err
		}
		consumer(val.Value.Value())
	}
	return nil
}

func (m *Metrics) UpdateAndGetMetrics(w http.ResponseWriter, r *http.Request) error {
	err := m.UpdateMemoryMetrics(); if err != nil {
		return err
	}
	err = m.UpdateRealtimeMetrics(); if err != nil {
		return err
	}
	promhttp.HandlerFor(m.registry, promhttp.HandlerOpts{}).ServeHTTP(w, r)
	return nil
}

func (m *Metrics) UpdateMemoryMetrics() error {
	file, err := os.Open("/proc/meminfo")
	if err != nil {
		return err
	}
	defer file.Close()

	var scan = bufio.NewScanner(file)

	for scan.Scan() {
		line := scan.Text()
		parts := strings.Fields(line)
		// Workaround for empty lines occasionally occur in CentOS 6.2 kernel 3.10.90.
		if len(parts) == 0 {
			continue
		}
		key := parts[0][:len(parts[0])-1] // remove trailing : from key
		// filter down to the fields of interest, skip parsing the rest
		gauge, ok := m.meminfoFields[key]; if !ok {
			continue
		}
		// Active(anon) -> Active_anon
		key = reParens.ReplaceAllString(key, "_${1}")
		fv, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return fmt.Errorf("invalid value in meminfo: %w", err)
		}
		switch len(parts) {
		case 2: // no unit
		case 3: // has unit, we presume kB
			fv *= 1024
			key = key + "_bytes"
		default:
			return fmt.Errorf("invalid line in meminfo: %s", line)
		}
		gauge.Set(fv)
	}

	return scan.Err()
}
