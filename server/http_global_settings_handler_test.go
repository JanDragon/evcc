package server

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/evcc-io/evcc/db/settings"
	"github.com/evcc-io/evcc/util/auth"
	"github.com/evcc-io/evcc/util/config"
	"github.com/stretchr/testify/assert"
)

func TestSettingsSetYamlHandlerCriticalPlugin(t *testing.T) {
	const pw = "secret"
	const key = "test_circuits_guard"
	body := "- name: main\n  getmaxcurrent:\n    source: script\n    cmd: echo 1"

	a := fakeAuth{mode: auth.Enabled, password: pw}
	h := settingsSetYamlHandler(key, []map[string]any{}, []config.Named{}, a)

	// session without password is rejected and nothing is persisted
	r := httptest.NewRequest(http.MethodPost, "/circuits", strings.NewReader(body))
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusPreconditionRequired, w.Code)
	_, err := settings.String(key)
	assert.ErrorIs(t, err, settings.ErrNotFound)

	// valid admin password persists the config
	r = httptest.NewRequest(http.MethodPost, "/circuits", strings.NewReader(body))
	r.Header.Set("X-Admin-Password", pw)
	w = httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	got, err := settings.String(key)
	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(body), got)
}

func TestSettingsSetYamlHandlerLevels(t *testing.T) {
	const key = "test_levels"
	body := "site: debug\ncache: error\n"

	h := settingsSetYamlHandler(key, map[string]any{}, map[string]string{}, fakeAuth{})
	r := httptest.NewRequest(http.MethodPost, "/levels", strings.NewReader(body))
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusOK, w.Code)
	got, err := settings.String(key)
	assert.NoError(t, err)
	assert.Equal(t, strings.TrimSpace(body), got)
}

func TestSettingsSetStringHandler(t *testing.T) {
	const key = "test_plant"

	var (
		publishedKey string
		publishedVal any
	)
	h := settingsSetStringHandler(key, func(key string, val any) {
		publishedKey = key
		publishedVal = val
	})

	r := httptest.NewRequest(http.MethodPost, "/plant", strings.NewReader("  foo-bar \n"))
	w := httptest.NewRecorder()
	h(w, r)

	assert.Equal(t, http.StatusOK, w.Code)

	got, err := settings.String(key)
	assert.NoError(t, err)
	assert.Equal(t, "foo-bar", got)
	assert.Equal(t, key, publishedKey)
	assert.Equal(t, "foo-bar", publishedVal)
}
