// nolint:errcheck // ignore cause tests, there won't be any errors
package api_test

import (
	"log"
	"os"
	"testing"

	"github.com/spf13/afero"
	"github.com/supabase/supabase-admin-api/api"

	"github.com/supabase/supabase-admin-api/test/data"
)

var (
	FS afero.Fs
	IO *afero.Afero
)

func TestMain(m *testing.M) {
	Logger := log.New(os.Stdout, "", 0)

	Logger.Println("Global tests setup")
	FS = afero.NewMemMapFs()
	IO = &afero.Afero{Fs: FS}
	api := api.NewAPIWithVersion(&api.Config{
		JwtSecret:                      "awdawdawdawdawdaw",
		UpstreamMetricsRefreshDuration: "60s",
		Port:                           8085,
		Host:                           "localhost",
	}, "0.0", FS)

	FS.MkdirAll("/etc", 0755)
	FS.MkdirAll("/etc/postgrest", 0755)
	FS.MkdirAll("/etc/kong", 0755)
	FS.MkdirAll("/etc/pgbouncer-custom", 0755)
	afero.WriteFile(FS, "/etc/postgrest/base.conf", []byte(data.PostgrestConf), 0644)
	afero.WriteFile(FS, "/etc/kong/kong.yml", []byte(data.KongConf), 0644)
	afero.WriteFile(FS, "/etc/pg_listen.conf", []byte(data.PglistenConf), 0644)
	afero.WriteFile(FS, "/etc/pgbouncer-custom/custom-overrides.conf", []byte(data.Pgbouncer), 0644)

	go func() {
		api.ListenAndServe("localhost:8085", "", "")
	}()

	// run tests
	exitVal := m.Run()

	Logger.Println("Global teardown")

	// exit process with tests exit code
	os.Exit(exitVal)
}
