package api

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/bluele/gcache"
	"github.com/go-chi/chi/middleware"
	"github.com/prometheus/common/expfmt"
	metrics "github.com/supabase/supabase-admin-api/api/metrics_endpoint"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/rs/cors"
	"github.com/sebest/xff"
	"github.com/sirupsen/logrus"
)

const (
	audHeaderName  = "X-JWT-AUD"
	defaultVersion = "unknown version"
)

// Config is the main API config
type Config struct {
	Host                           string                        `yaml:"host" default:"localhost"`
	Port                           int                           `yaml:"port" default:"8085"`
	JwtSecret                      string                        `yaml:"jwt_secret" required:"true"`
	MetricCollectors               []string                      `yaml:"metric_collectors" required:"true"`
	GotrueHealthEndpoint           string                        `yaml:"gotrue_health_endpoint" required:"false" default:"http://localhost:9999/health"`
	PostgrestEndpoint              string                        `yaml:"postgrest_endpoint" required:"false" default:"http://localhost:3000/"`
	RealtimeServiceName            string                        `yaml:"realtime_service_name" required:"false" default:"supabase"`
	UpstreamMetricsSources         []metrics.MetricsSourceConfig `yaml:"upstream_metrics_sources" required:"true"`
	NodeExporterAdditionalArgs     []string                      `yaml:"node_exporter_additional_args" required:"false"`
	UpstreamMetricsRefreshDuration string                        `yaml:"upstream_metrics_refresh_duration" default:"60s"`

	// supply to enable TLS termination
	KeyPath  string `yaml:"key_path" required:"false"`
	CertPath string `yaml:"cert_path" required:"false"`
}

func (c *Config) GetMetricsSources() []metrics.MetricsSource {
	logger := logrus.New()
	var parser expfmt.TextParser
	sources := make([]metrics.MetricsSource, 0)
	for _, config := range c.UpstreamMetricsSources {
		client := http.Client{
			Timeout: 1 * time.Second,
			Transport: &http.Transport{
				TLSClientConfig: &tls.Config{
					InsecureSkipVerify: config.SkipTlsVerify,
				},
			},
		}
		sourceLogger := logger.WithField("source", config.Name)
		source := metrics.MetricsSource{
			Config:     config,
			HttpClient: &client,
			Logger:     sourceLogger,
			Parser:     &parser,
		}
		logger.Infof("Creating source for %+v", config)
		sources = append(sources, source)
	}
	return sources
}

// API is the main REST API
type API struct {
	handler        http.Handler
	config         *Config
	version        string
	metricsSources []metrics.MetricsSource
}

// ListenAndServe starts the REST API
func (a *API) ListenAndServe(hostAndPort string, keyPath string, certPath string) {
	log := logrus.WithField("component", "api")
	server := &http.Server{
		Addr:    hostAndPort,
		Handler: a.handler,
	}

	done := make(chan struct{})
	defer close(done)
	go func() {
		waitForTermination(log, done)
		ctx, cancel := context.WithTimeout(context.Background(), time.Minute)
		defer cancel()
		server.Shutdown(ctx)
	}()

	var err error
	if keyPath != "" && certPath != "" {
		log.WithField("cert", certPath).WithField("key", keyPath).Info("Using TLS")
		err = server.ListenAndServeTLS(certPath, keyPath)
	} else {
		log.Warn("Not using TLS!")
		err = server.ListenAndServe()
	}
	if err != http.ErrServerClosed {
		log.WithError(err).Fatal("http server listen failed")
	}
}

// WaitForShutdown blocks until the system signals termination or done has a value
func waitForTermination(log logrus.FieldLogger, done <-chan struct{}) {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)
	select {
	case sig := <-signals:
		log.Infof("Triggering shutdown from signal %s", sig)
	case <-done:
		log.Infof("Shutting down...")
	}
}

// NewAPIWithVersion creates a new REST API using the specified version
func NewAPIWithVersion(config *Config, version string) *API {
	api := &API{config: config, version: version}
	nodeMetrics, err := NewMetrics(config.MetricCollectors, config.GotrueHealthEndpoint, config.PostgrestEndpoint, config.NodeExporterAdditionalArgs)
	if err != nil {
		panic(fmt.Sprintf("Couldn't initialize metrics: %+v", err))
	}

	projectMetrics := metrics.Metrics{
		Sources: config.GetMetricsSources(),
	}
	duration, err := time.ParseDuration(config.UpstreamMetricsRefreshDuration)
	if err != nil {
		logrus.WithError(err).Fatal("failed to parse metrics refresh duration")
	}
	cache := gcache.New(1).Expiration(duration).LoaderFunc(func(_ interface{}) (interface{}, error) {
		return projectMetrics.GetMergedMetrics(), nil
	}).Build()
	xffmw, _ := xff.Default()

	r := chi.NewRouter()
	r.Use(xffmw.Handler)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// unauthenticated
	r.Group(func(r chi.Router) {
		r.Method("GET", "/metrics", nodeMetrics.GetHandler())
		r.Method("GET", "/project-metrics", ErrorHandlingWrapper(api.ServeUpstreamMetrics(cache.Get)))
	})

	// private endpoints
	r.Group(func(r chi.Router) {
		r.Use(api.AuthHandler)
		r.Method("GET", "/health", ErrorHandlingWrapper(api.HealthCheck))

		r.Route("/", func(r chi.Router) {
			r.Route("/test", func(r chi.Router) {
				r.Method("GET", "/", ErrorHandlingWrapper(api.TestGet))
			})

			r.Route("/service", func(r chi.Router) {
				// applications are kong, pglisten, postgrest, goauth, realtime, adminapi, all
				r.Route("/restart", func(r chi.Router) {
					r.Method("GET", "/", ErrorHandlingWrapper(api.RestartServices))
					r.Method("GET", "/{application}", ErrorHandlingWrapper(api.RestartServices))
				})
				r.Method("GET", "/reboot", ErrorHandlingWrapper(api.RebootMachine))
			})

			// applications are kong, pglisten, postgrest, goauth, realtime, adminapi
			r.Route("/config/{application}", func(r chi.Router) {
				r.Method("GET", "/", ErrorHandlingWrapper(api.GetFileContents))
				r.Method("POST", "/", ErrorHandlingWrapper(api.SetFileContents))
			})

			// applications are kong, pglisten, postgrest, goauth, realtime
			r.Route("/logs/{application}/{type}/{n:[0-9]*}", func(r chi.Router) {
				r.Method("GET", "/", ErrorHandlingWrapper(api.GetLogContents))
			})

			r.Route("/cert", func(r chi.Router) {
				r.Method("POST", "/", ErrorHandlingWrapper(api.UpdateCert))
			})

			r.Route("/disk", func(r chi.Router) {
				r.Method("POST", "/expand", ErrorHandlingWrapper(ExpandFilesystem))
			})
		})
	})

	corsHandler := cors.New(cors.Options{
		AllowedMethods:   []string{http.MethodGet, http.MethodPost, http.MethodPut, http.MethodDelete},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", audHeaderName},
		AllowCredentials: true,
	})

	api.handler = corsHandler.Handler(chi.ServerBaseContext(context.Background(), r))
	return api
}

// HealthCheck returns basic information for status purposes
func (a *API) HealthCheck(w http.ResponseWriter, r *http.Request) error {
	return sendJSON(w, http.StatusOK, map[string]string{
		"version":     a.version,
		"name":        "supabase-admin-api",
		"description": "supabase-admin-api is an api to manage KPS",
	})
}
