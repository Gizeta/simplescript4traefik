package simplescript4traefik

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Gizeta/simplescript4traefik"
)

func TestDemo(t *testing.T) {
	cfg := simplescript4traefik.CreateConfig()
	cfg.Code = `
		(set_req_header "X-Demo" "test")
	`

	ctx := context.Background()
	next := http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {})

	handler, err := simplescript4traefik.New(ctx, next, cfg, "demo-plugin")
	if err != nil {
		t.Fatal(err)
	}

	recorder := httptest.NewRecorder()

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, "http://localhost", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(recorder, req)

	assertHeader(t, req, "X-Demo", "test")
}

func assertHeader(t *testing.T, req *http.Request, key, expected string) {
	t.Helper()

	if req.Header.Get(key) != expected {
		t.Errorf("invalid header value: %s", req.Header.Get(key))
	}
}
