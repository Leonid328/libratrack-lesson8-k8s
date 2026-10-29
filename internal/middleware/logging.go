package middleware

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

func Logging() gin.HandlerFunc {
	return func(c *gin.Context) {
		var body []byte
		if c.Request.Body != nil {
			body, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(body))
		}
		safe := redact(body)
		log.Printf("%s %s %s", c.Request.Method, c.Request.URL.Path, safe)
		c.Next()
	}
}

func redact(body []byte) string {
	if len(body) == 0 {
		return ""
	}
	var m map[string]any
	if err := json.Unmarshal(body, &m); err != nil {
		return ""
	}
	for k := range m {
		if strings.EqualFold(k, "password") {
			m[k] = "[redacted]"
		}
	}
	b, _ := json.Marshal(m)
	return string(b)
}
