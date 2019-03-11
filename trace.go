// Package echotrace enables golang.org/x/net/trace based request tracing for
// echo based servers.
package echotrace

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"golang.org/x/net/trace"
)

const (
	// ContextKey is the key used to lookup the trace.Trace for the current
	// request from the echo.Context.
	ContextKey = "trace-context-key"
)

// Middleware is a echo.MiddlewareFunc that creates a new trace.Trace for the
// current request and sets in on the echo.Context.
func Middleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		tr := trace.New(c.Request().Method+" "+c.Request().URL.Path, c.Request().URL.String())
		tr.LazyPrintf("From: %s", c.Request().RemoteAddr)
		defer tr.Finish()
		c.Set(ContextKey, tr)
		err := next(c)
		if err != nil {
			tr.LazyPrintf("Error: [%v]", err)
			tr.SetError()
		} else if status := c.Response().Status; status != http.StatusOK {
			tr.LazyPrintf("Error: Code=%d", status)
			tr.SetError()
		}
		return err
	}
}

// Handler is a echo.HandlerFunc that serves the page of traces.
var Handler = echo.WrapHandler(http.HandlerFunc(trace.Traces))
