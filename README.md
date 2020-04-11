# elton-etag

The middleware has been archived, please use the middleware of [elton](https://github.com/vicanso/elton).

[![Build Status](https://img.shields.io/travis/vicanso/elton-etag.svg?label=linux+build)](https://travis-ci.org/vicanso/elton-etag)

ETag middleware for elton, generate response's ETag header by `sha1`.

```go
package main

import (
	"bytes"

	"github.com/vicanso/elton"

	etag "github.com/vicanso/elton-etag"
)

func main() {

	e := elton.New()
	e.Use(etag.NewDefault())

	e.GET("/", func(c *elton.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("abcd")
		return
	})

	err := e.ListenAndServe(":3000")
	if err != nil {
		panic(err)
	})
}

```

