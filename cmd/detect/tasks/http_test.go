package tasks

import (
	"fmt"
	"net"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHttpHandlerQueryFiltersResponses(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{name: "match", body: `{"match":true}`, want: true},
		{name: "no match", body: `{"other":true}`, want: false},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				fmt.Fprint(w, tc.body)
			}))
			defer srv.Close()

			host, port, err := net.SplitHostPort(srv.Listener.Addr().String())
			require.NoError(t, err)

			handler, err := HttpHandlerFactory(map[string]any{
				"schema":  "http",
				"jq":      ".match",
				"timeout": 500 * time.Millisecond,
			})
			require.NoError(t, err)

			got := handler.Test(nil, ResultDetails{IP: host, Port: mustAtoi(t, port)})
			if tc.want {
				require.Len(t, got, 1)
				assert.Equal(t, host, got[0].IP)
				assert.Equal(t, mustAtoi(t, port), got[0].Port)
			} else {
				assert.Nil(t, got)
			}
		})
	}
}

func mustAtoi(t *testing.T, s string) int {
	t.Helper()

	port, err := net.LookupPort("tcp", s)
	require.NoError(t, err)

	return port
}
