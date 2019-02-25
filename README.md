# cod-etag

[![Build Status](https://img.shields.io/travis/vicanso/cod-etag.svg?label=linux+build)](https://travis-ci.org/vicanso/cod-etag)

ETag middleware for cod, generate response's ETag header by `sha1`.

```go
package main

import (
	"bytes"

	"github.com/vicanso/cod"

	etag "github.com/vicanso/cod-etag"
)

func main() {

	d := cod.New()
	d.Use(etag.NewDefault())

	d.GET("/", func(c *cod.Context) (err error) {
		c.BodyBuffer = bytes.NewBufferString("abcd")
		return
	})

	d.ListenAndServe(":7001")
}

```

