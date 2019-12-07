// Copyright 2018 tree xie
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package etag

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/vicanso/elton"
)

type (
	// Config eTag config
	Config struct {
		Skipper elton.Skipper
	}
)

// gen generate eTag
func gen(buf []byte) string {
	size := len(buf)
	if size == 0 {
		return `"0-2jmj7l5rSw0yVb_vlWAYkK_YBwk="`
	}
	h := sha1.New()
	_, err := h.Write(buf)
	if err != nil {
		return ""
	}
	hash := base64.URLEncoding.EncodeToString(h.Sum(nil))
	return fmt.Sprintf(`"%x-%s"`, size, hash)
}

// NewDefault create a default ETag middleware
func NewDefault() elton.Handler {
	return New(Config{})
}

// New create a ETag middleware
func New(config Config) elton.Handler {
	skipper := config.Skipper
	if skipper == nil {
		skipper = elton.DefaultSkipper
	}
	return func(c *elton.Context) (err error) {
		if skipper(c) {
			return c.Next()
		}
		err = c.Next()
		if err != nil {
			return
		}
		bodyBuf := c.BodyBuffer
		// 如果无内容或已设置 ETag ，则跳过
		// 因为没有内容也不生成 ETag
		if bodyBuf == nil || bodyBuf.Len() == 0 ||
			c.GetHeader(elton.HeaderETag) != "" {
			return
		}
		// 如果响应状态码不为0 而且( < 200 或者 >= 300)，则跳过
		// 如果未设置状态码，最终为200
		statusCode := c.StatusCode
		if statusCode != 0 &&
			(statusCode < http.StatusOK ||
				statusCode >= http.StatusMultipleChoices) {
			return
		}
		eTag := gen(bodyBuf.Bytes())
		if eTag != "" {
			c.SetHeader(elton.HeaderETag, eTag)
		}
		return
	}
}
