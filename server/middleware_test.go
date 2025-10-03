package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mahmoud-shabban/snippetbox/internal/assert"
)

func TestCommonHeaders(t *testing.T) {

	rec := httptest.NewRecorder()

	req, err := http.NewRequest(http.MethodGet, "/", nil)

	if err != nil {
		t.Fatal(err)
	}

	next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("ok!"))
	})

	commonHeaders(next).ServeHTTP(rec, req)

	result := rec.Result()

	securityHeader := result.Header.Get("Content-Security-Policy")
	securityHeaderExpected := "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com"

	assert.Equal(t, securityHeader, securityHeaderExpected)

	reffererPolicyHeader := result.Header.Get("Referrer-Policy")
	reffererPolicyHeaderExpected := "origin-when-cross-origin"

	assert.Equal(t, reffererPolicyHeader, reffererPolicyHeaderExpected)

	contentTypeHeader := result.Header.Get("X-Content-Type-Options")
	contentTypeHeaderExpected := "nosniff"

	assert.Equal(t, contentTypeHeader, contentTypeHeaderExpected)

	frameOptionsHeader := result.Header.Get("X-Frame-Options")
	frameOptionsHeaderExpected := "deny"

	assert.Equal(t, frameOptionsHeader, frameOptionsHeaderExpected)

	xssHeader := result.Header.Get("X-XSS-Protection")
	xssHeaderExpected := "0"

	assert.Equal(t, xssHeader, xssHeaderExpected)

	serverHeader := result.Header.Get("Server")
	serverHeaderExpected := "Go"

	assert.Equal(t, serverHeader, serverHeaderExpected)

}
