package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/vuquang23/poseidon/internal/pkg/util/requestid"
	"github.com/vuquang23/poseidon/pkg/logger"
)

func NewLoggerMiddleware(logCfg logger.Config, logBackend logger.LoggerBackend) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		requestID := requestid.ExtractRequestID(c)

		var buf bytes.Buffer
		tee := io.TeeReader(c.Request.Body, &buf)
		reqBody, _ := io.ReadAll(tee)
		c.Request.Body = io.NopCloser(&buf)

		commonFields := logger.Fields{
			"request.id": requestID,
		}

		reqLogger := logger.WithFieldsNonContext(commonFields)
		c.Set(string(logger.CtxLoggerKey), reqLogger)

		reqLogger.WithFields(logger.Fields{
			"request.method":     c.Request.Method,
			"request.uri":        c.Request.URL.RequestURI(),
			"request.body":       string(reqBody),
			"request.client_ip":  c.ClientIP(),
			"request.user_agent": c.Request.UserAgent(),
		}).Info("inbound request")

		blw := &bodyLogWriter{body: bytes.NewBufferString(""), ResponseWriter: c.Writer}
		c.Writer = blw

		c.Next()

		resp, _ := io.ReadAll(blw.body)

		reqLogger.WithFields(
			logger.Fields{
				"response.status":      blw.Status(),
				"response.body":        string(resp),
				"response.duration_ms": time.Since(startTime).Milliseconds(),
			}).
			Info("inbound response")
	}
}

type bodyLogWriter struct {
	gin.ResponseWriter
	body *bytes.Buffer
}

func (w bodyLogWriter) Write(b []byte) (int, error) {
	w.body.Write(b)
	return w.ResponseWriter.Write(b)
}

func (w bodyLogWriter) WriteString(s string) (int, error) {
	w.body.WriteString(s)
	return w.ResponseWriter.WriteString(s)
}
