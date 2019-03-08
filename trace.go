package echotrace

import (
	"net/http"

	"github.com/labstack/echo"
	"golang.org/x/net/trace"
)

const (
	ContextKey = "trace-context-key"
)

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
		}
		return err
	}
}

var Handler = echo.WrapHandler(http.HandlerFunc(trace.Traces))
