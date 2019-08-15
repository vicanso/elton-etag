# elton-etag

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

	d := elton.New()
	d.Use(etag.NewDefault())

	d.GET("/", func(c *elton.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("abcd")
		return
	})

	d.ListenAndServe(":7001")
}

```

