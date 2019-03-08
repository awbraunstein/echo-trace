# echo-trace
Request tracing for the the echo framework based on golang.org/x/net/trace

# Installation

`go get -u github.com/awbraunstein/echo-trace`

# Usage

```go
package main

import (
	echotrace "github.com/awbraunstein/echo-trace"
	"github.com/labstack/echo"
	"golang.org/x/net/trace"
)

func main() {
	e := echo.New()
	// Install the echotrace middleware.
	e.Use(echotrace.Middleware)

	// Enable the echotrace handler on /debug/requests.
	e.GET("/debug/requests", echotrace.Handler)

	e.GET("/", func(c echo.Context) error {
		tr := c.Get(ehcotrace.ContextKey).(trace.Trace)
		tr.LazyPrintf("Adding some data to the trace")
		return c.String(http.StatusOK, "Handled /")
	})
}
```
