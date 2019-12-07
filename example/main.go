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
	}
}
