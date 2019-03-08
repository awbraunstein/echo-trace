# echo-trace
[![MIT licensed](https://img.shields.io/github/license/awbraunstein/echo-trace.svg)](LICENSE)
[![GoDoc](https://godoc.org/github.com/awbraunstein/echo-trace?status.svg)](https://godoc.org/github.com/awbraunstein/echo-trace)
[![Build Status](https://travis-ci.com/awbraunstein/echo-trace.svg?branch=master)](https://travis-ci.com/awbraunstein/echo-trace)
[![Coverage Status](https://img.shields.io/codecov/c/github/awbraunstein/echo-trace.svg)](https://codecov.io/gh/awbraunstein/echo-trace)
[![Go Report Card](https://goreportcard.com/badge/github.com/awbraunstein/echo-trace)](https://goreportcard.com/report/github.com/awbraunstein/echo-trace)

Request tracing for the the echo framework based on [golang.org/x/net/trace](https://godoc.org/golang.org/x/net/trace).

# Installation

`go get -u github.com/awbraunstein/echo-trace`

# Usage

```go
package main

import (
	"net/http"

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
