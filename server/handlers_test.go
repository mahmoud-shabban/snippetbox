package main

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mahmoud-shabban/snippetbox/internal/assert"
)

func TestPing(t *testing.T) {
	rec := httptest.NewRecorder()

	req, err := http.NewRequest("GET", "/healthz", nil)
	if err != nil {
		t.Fatal(err)
	}

	ping(rec, req)

	result := rec.Result()

	assert.Equal(t, result.StatusCode, http.StatusOK)

	buf, err := io.ReadAll(result.Body)

	if err != nil {
		t.Fatal(err)
	}

	defer result.Body.Close()

	body := bytes.TrimSpace(buf)
	assert.Equal(t, string(body), "ok!")

}
