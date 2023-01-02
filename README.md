# gaglog [![go.dev reference](https://img.shields.io/badge/go.dev-reference-007d9c?logo=go&logoColor=white&style=flat-square)](https://pkg.go.dev/github.com/kyslik/gaglog/v1) [![Latest Release](https://img.shields.io/github/v/release/kyslik/gaglog?style=flat-square)](https://github.com/kyslik/gaglog/releases/latest) ![Build Status](https://github.com/kyslik/gaglog/actions/workflows/test.yaml/badge.svg?branch=main)

`gaglog` is a Go package that adds regex/duration based filtering to the standard library [`log`](https://pkg.go.dev/log) package.

This package was heavily influenced by [`hashicorp/logutils`](https://github.com/hashicorp/logutils) package.

## Simple example

Presumably your application already uses the default `log` package. To start throttling logs, you'll want your code to look like the following:

```go
package main

import (
	"log"
	"os"
	"regexp"
	"time"

	"github.com/kyslik/gaglog"
)

func main() {
	filter := &gaglog.GagFilter{
		Writer: os.Stderr,
		Gags: gaglog.Gags{
			regexp.MustCompile("P([a-z]+)ch"): time.Millisecond * 100,
			regexp.MustCompile("L([a-z]+)ch"): time.Millisecond * 1000,
			regexp.MustCompile("F([a-z]+)ch"): time.Millisecond * 10000,
		},
	}

	log.SetOutput(filter)

	for i := 0; i < 1000000; i++ {
		log.Print("Pinch")
		log.Print("Lynch")
		log.Print("Flinch")
		if i % 100000 == 0 {
			log.Print("Grinch")
		}
	}
	time.Sleep(5 * time.Second)
	log.Println("Clinch")
}
```

This logs to standard error exactly like Go's standard logger. Any log messages that are not matched by a regex, won't be gaged.
