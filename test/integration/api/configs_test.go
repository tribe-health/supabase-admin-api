package api_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/supabase/supabase-admin-api/api"
	"github.com/supabase/supabase-admin-api/test/common"
	"github.com/supabase/supabase-admin-api/test/data"
)

func TestGetConf(t *testing.T) {
	cases := []struct {
		description string
		url         string
		expected    string
	}{
		{
			description: "postgrest",
			url:         "http://localhost:8085/config/postgrest",
			expected:    data.PostgrestConf,
		},
		{
			description: "kong",
			url:         "http://localhost:8085/config/kong",
			expected:    data.KongConf,
		},
		{
			description: "pglisten",
			url:         "http://localhost:8085/config/pglisten",
			expected:    data.PglistenConf,
		},
	}

	for i := range cases {
		tt := cases[i]
		t.Run(tt.description, func(t *testing.T) {
			req, err := http.NewRequest("GET", tt.url, nil)
			require.NoError(t, err)
			respBody, err := common.AuthedRequest(req)
			require.NoError(t, err)

			fc := api.FileContents{
				RawContents:     tt.expected,
				RestartServices: false,
			}
			expectedJson, _ := json.Marshal(fc)
			require.Equal(t, string(expectedJson), string(respBody))
		})
	}
}

func TestPostConf(t *testing.T) {
	cases := []struct {
		description string
		url         string
		data        string
		expected    string
	}{
		{
			description: "pgbouncer",
			url:         "http://localhost:8085/config/pgbouncer",
			data:        data.PgbouncerNew,
			expected:    "/etc/pgbouncer-custom/custom-overrides.conf",
		},
	}

	for i := range cases {
		tt := cases[i]
		t.Run(tt.description, func(t *testing.T) {
			fc := api.FileContents{
				RawContents:     tt.data,
				RestartServices: false,
			}
			body, _ := json.Marshal(fc)
			req, err := http.NewRequest("POST", tt.url, bytes.NewBuffer(body))
			require.NoError(t, err)
			_, err = common.AuthedRequest(req)
			require.NoError(t, err)

			actual, err := IO.ReadFile(tt.expected)
			require.NoError(t, err)
			require.Equal(t, tt.data, string(actual))
		})
	}
}
