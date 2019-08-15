package etag

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"errors"
	"fmt"
	"math/rand"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vicanso/elton"
)

var testData []byte

func init() {
	var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

	fn := func(n int) string {
		b := make([]rune, n)
		for i := range b {
			b[i] = letterRunes[rand.Intn(len(letterRunes))]
		}
		return string(b)
	}
	testData = []byte(fn(4096))
}

func TestGen(t *testing.T) {
	assert := assert.New(t)
	assert.Equal(gen([]byte("")), `"0-2jmj7l5rSw0yVb_vlWAYkK_YBwk="`)
}

func TestETag(t *testing.T) {
	fn := NewDefault()
	t.Run("curstom error", func(t *testing.T) {
		assert := assert.New(t)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, nil)
		customErr := errors.New("abcd")
		c.Next = func() error {
			return customErr
		}
		err := fn(c)
		assert.Equal(err, customErr)
	})

	t.Run("no body", func(t *testing.T) {
		assert := assert.New(t)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, nil)
		c.Next = func() error {
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.Empty(c.GetHeader(elton.HeaderETag))
	})

	t.Run("error status", func(t *testing.T) {
		assert := assert.New(t)
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, nil)
		c.Next = func() error {
			c.Body = map[string]string{
				"name": "tree.xie",
			}
			c.StatusCode = 400
			c.BodyBuffer = bytes.NewBufferString(`{"name":"tree.xie"}`)
			return nil
		}
		err := fn(c)
		assert.Nil(err)
		assert.Empty(c.GetHeader(elton.HeaderETag))
	})

	t.Run("gen eTag", func(t *testing.T) {
		resp := httptest.NewRecorder()
		c := elton.NewContext(resp, nil)
		c.Next = func() error {
			c.Body = map[string]string{
				"name": "tree.xie",
			}
			c.BodyBuffer = bytes.NewBufferString(`{"name":"tree.xie"}`)
			return nil
		}
		err := fn(c)
		if err != nil {
			t.Fatalf("eTag middleware fail, %v", err)
		}
		if c.GetHeader(elton.HeaderETag) != `"13-yo9YroUOjW1obRvVoXfrCiL2JGE="` {
			t.Fatalf("gen eTag fail")
		}
	})
}

func BenchmarkGenETag(b *testing.B) {
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		gen(testData)
	}
}

func BenchmarkMd5(b *testing.B) {
	b.ReportAllocs()
	fn := func(buf []byte) string {
		size := len(buf)
		if size == 0 {
			return `"0-2jmj7l5rSw0yVb_vlWAYkK_YBwk="`
		}
		h := md5.New()
		h.Write(buf)
		hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
		return fmt.Sprintf(`"%x-%s"`, size, hash)
	}
	for i := 0; i < b.N; i++ {
		fn(testData)
	}
}

// https://stackoverflow.com/questions/50120427/fail-unit-tests-if-coverage-is-below-certain-percentage
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	rc := m.Run()

	// rc 0 means we've passed,
	// and CoverMode will be non empty if run with -cover
	if rc == 0 && testing.CoverMode() != "" {
		c := testing.Coverage()
		if c < 0.9 {
			fmt.Println("Tests passed but coverage failed at", c)
			rc = -1
		}
	}
	os.Exit(rc)
}
