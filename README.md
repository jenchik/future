# future
Golang implementation of futures/promises. Idea from https://github.com/sentientmonkey/future

[![Build Status](https://travis-ci.org/jenchik/future.svg)](https://travis-ci.org/jenchik/future)
[![GoDoc](https://godoc.org/github.com/jenchik/future?status.svg)](https://godoc.org/github.com/jenchik/future)
[![Go Report Card](https://goreportcard.com/badge/github.com/jenchik/future)](https://goreportcard.com/report/github.com/jenchik/future)

Installation
------------

```bash
go get github.com/jenchik/future
```

Example
-------
```go
package main

import (
	"fmt"
	"net/http"

	"github.com/jenchik/future"
)

func main() {
	f := future.NewFuture(func() (interface{}, error) {
		return http.Get("http://golang.org/")
	})

	result, err := f.Wait(nil)
	if err != nil {
		fmt.Printf("Got error: %s\n", err)
		return
	}

	response := result.(*http.Response)
	defer response.Body.Close()

	fmt.Printf("Got result: %d\n", response.StatusCode)
}
```
