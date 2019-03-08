package echotrace

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo"
	"golang.org/x/net/trace"
)

func TestMiddleware(t *testing.T) {
	e := echo.New()
	e.Use(Middleware)
	// For the test, we need to disable the auth function that protects the
	// traces.
	trace.AuthRequest = func(req *http.Request) (any, sensitive bool) {
		return true, true
	}
	e.GET("/debug/requests", Handler)
	e.GET("/", func(c echo.Context) error {
		tr := c.Get(ContextKey).(trace.Trace)
		tr.LazyPrintf("Adding some data to the trace")
		return c.String(http.StatusOK, "DONE")
	})
	e.GET("/error", func(c echo.Context) error {
		tr := c.Get(ContextKey).(trace.Trace)
		tr.LazyPrintf("Adding some data to the trace")
		return c.String(http.StatusBadRequest, "ERROR")
	})
	// Tests good request.
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected OK status, but got %d", rec.Code)
	}
	if got := rec.Body.String(); got != "DONE" {
		t.Fatalf("Expected DONE, but got: %s", got)
	}

	// Tests error request.
	req = httptest.NewRequest(http.MethodGet, "/error", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusBadRequest {
		t.Fatalf("Expected 400 status, but got %d", rec.Code)
	}
	if got := rec.Body.String(); got != "ERROR" {
		t.Fatalf("Expected ERROR, but got: %s", got)
	}

	// Tests invalid url request.
	req = httptest.NewRequest(http.MethodGet, "/invalid", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusNotFound {
		t.Fatalf("Expected 404 status, but got %d", rec.Code)
	}
	want := "{\"message\":\"Not Found\"}\n"
	if got := rec.Body.String(); got != want {
		t.Fatalf("Expected %s, but got: %s", want, got)
	}

	// Tests debug/requests handler is installed.
	req = httptest.NewRequest(http.MethodGet, "/debug/requests", nil)
	rec = httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("Expected OK status, but got %d", rec.Code)
	}
	if got := rec.Body.String(); !strings.Contains(got, "/error") {
		t.Fatalf("Expected response to contain /error, but was: %s", got)
	}
}
