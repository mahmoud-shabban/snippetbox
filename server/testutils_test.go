package main

import (
	"bytes"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-playground/form/v4"
	"github.com/mahmoud-shabban/snippetbox/internal/models/mocks"
)

type testServer struct {
	*httptest.Server
}

func newTestServer(t *testing.T, handler http.Handler) *testServer {

	ts := httptest.NewTLSServer(handler)

	jar, err := cookiejar.New(nil)
	if err != nil {
		t.Fatal(err)
	}

	ts.Client().Jar = jar

	ts.Client().CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}

	return &testServer{ts}
}

func newTestApp(t *testing.T) *Application {
	sessionManager := scs.New()
	sessionManager.Lifetime = 12 * time.Hour
	sessionManager.Cookie.Secure = true

	tempCache, err := newTemplateCache()
	if err != nil {
		t.Fatal(err)
	}

	formDecoder := form.NewDecoder()
	app := Application{
		logger:         slog.New(slog.DiscardHandler),
		sessionManager: sessionManager,
		templateCache:  tempCache,
		formDecoder:    formDecoder,
		users:          &mocks.UserModel{},
		snippets:       &mocks.SnippetModel{},
	}
	return &app
}

func (srv *testServer) get(t *testing.T, endpoint string) (int, http.Header, string) {
	res, err := srv.Client().Get(srv.URL + endpoint)

	if err != nil {
		t.Fatal(err)
	}

	buf, err := io.ReadAll(res.Body)

	if err != nil {
		t.Fatal(err)
	}

	buf = bytes.TrimSpace(buf)

	return res.StatusCode, res.Header, string(buf)
}
