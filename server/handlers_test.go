package main

import (
	"net/http"
	"testing"

	"github.com/mahmoud-shabban/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	app := newTestApp(t)

	ts := newTestServer(t, app.routes())
	defer ts.Close()

	status, _, body := ts.get(t, "/healthz")

	assert.Equal(t, status, http.StatusOK)

	assert.Equal(t, body, "ok!")

}

// func TestPing(t *testing.T) {
// 	rec := httptest.NewRecorder()

// 	req, err := http.NewRequest("GET", "/healthz", nil)
// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	ping(rec, req)

// 	result := rec.Result()

// 	assert.Equal(t, result.StatusCode, http.StatusOK)

// 	buf, err := io.ReadAll(result.Body)

// 	if err != nil {
// 		t.Fatal(err)
// 	}

// 	defer result.Body.Close()

// 	body := bytes.TrimSpace(buf)
// 	assert.Equal(t, string(body), "ok!")

// }
