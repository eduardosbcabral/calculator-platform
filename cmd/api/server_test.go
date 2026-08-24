package main

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestHandlerWithoutFrontend(t *testing.T) {
	handler, err := newHandler("")
	if err != nil {
		t.Fatal(err)
	}

	tests := []struct {
		name       string
		method     string
		path       string
		wantStatus int
		wantBody   string
	}{
		{name: "health", method: http.MethodGet, path: "/healthz", wantStatus: http.StatusOK, wantBody: `{"status":"ok"}`},
		{name: "health method", method: http.MethodPost, path: "/healthz", wantStatus: http.StatusMethodNotAllowed},
		{name: "frontend unavailable", method: http.MethodGet, path: "/", wantStatus: http.StatusNotFound},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			recorder := httptest.NewRecorder()
			handler.ServeHTTP(recorder, httptest.NewRequest(test.method, test.path, nil))

			if recorder.Code != test.wantStatus {
				t.Fatalf("status = %d, want %d", recorder.Code, test.wantStatus)
			}
			if test.wantBody != "" && strings.TrimSpace(recorder.Body.String()) != test.wantBody {
				t.Fatalf("body = %q, want %q", recorder.Body.String(), test.wantBody)
			}
		})
	}
}

func TestHandlerServesFrontend(t *testing.T) {
	staticDir := t.TempDir()
	if err := os.WriteFile(filepath.Join(staticDir, "index.html"), []byte("<h1>Calculator</h1>"), 0o600); err != nil {
		t.Fatal(err)
	}

	handler, err := newHandler(staticDir)
	if err != nil {
		t.Fatal(err)
	}
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(http.MethodGet, "/", nil))

	if recorder.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", recorder.Code, http.StatusOK)
	}
	if body := recorder.Body.String(); body != "<h1>Calculator</h1>" {
		t.Fatalf("body = %q", body)
	}
}

func TestHandlerRequiresFrontendEntryPoint(t *testing.T) {
	if _, err := newHandler(t.TempDir()); err == nil {
		t.Fatal("newHandler() error = nil, want an error")
	}
}

func TestNewServer(t *testing.T) {
	server := newServer(http.NewServeMux(), "9090")
	if server.Addr != ":9090" {
		t.Fatalf("address = %q, want :9090", server.Addr)
	}
}
