package api

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"regexp"
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

var bearerRegexp = regexp.MustCompile(`^(?:B|b)earer (\S+$)`)

// Config is the main API config
type Config struct {
	Host            string
	Port            int `envconfig:"PORT" default:"8085"`
	Endpoint        string
	RequestIDHeader string `envconfig:"REQUEST_ID_HEADER"`
	ExternalURL     string `json:"external_url" envconfig:"API_EXTERNAL_URL"`
}

// API is the main REST API
type API struct {
	handler http.Handler
	config  *Config
	version string
}

// ListenAndServe starts the REST API
func (a *API) ListenAndServe(hostAndPort string) {
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

	if err := server.ListenAndServe(); err != http.ErrServerClosed {
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

// NewAPI instantiates a new REST API
func NewAPI(config *Config) *API {
	return NewAPIWithVersion(config, defaultVersion)
}

// NewAPIWithVersion creates a new REST API using the specified version
func NewAPIWithVersion(config *Config, version string) *API {
	api := &API{config: config, version: version}

	xffmw, _ := xff.Default()

	r := newRouter()
	r.UseBypass(xffmw.Handler)
	r.Use(recoverer)

	r.Get("/health", api.HealthCheck)

	r.Route("/", func(r *router) {
		r.Route("/test", func(r *router) {
			r.Get("/", api.TestGet)
		})

		r.Route("/service", func(r *router) {
			// applications are kong, pglisten, postgrest, goauth, realtime, adminapi, all
			r.Route("/restart", func(r *router) {
				r.Get("/", api.RestartServices)
				r.Get("/{application}", api.RestartServices)
			})
			r.Get("/reboot", api.RebootMachine)
		})

		// applications are kong, pglisten, postgrest, goauth, realtime, adminapi
		r.Route("/config/{application}", func(r *router) {
			r.Get("/", api.GetFileContents)
			r.Post("/", api.SetFileContents)
		})

		// applications are kong, pglisten, postgrest, goauth, realtime
		r.Route("/logs/{application}/{type}/{n:[0-9]*}", func(r *router) {
			r.Get("/", api.GetLogContents)
		})

		r.Route("/cert", func(r *router) {
			r.Post("/", api.UpdateCert)
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
