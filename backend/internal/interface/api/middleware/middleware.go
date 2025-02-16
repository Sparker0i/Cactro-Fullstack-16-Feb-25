package middleware

import (
	"bytes"
	"io"
	"time"

	"github.com/Sparker0i/cactro-polls/internal/infrastructure/logger"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type Middleware struct {
	logger logger.Logger
}

func NewMiddleware(logger logger.Logger) *Middleware {
	return &Middleware{
		logger: logger,
	}
}

// RequestID adds a unique request ID to each request
func (m *Middleware) RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")
		if requestID == "" {
			requestID = uuid.New().String()
		}
		c.Set("request_id", requestID)
		c.Header("X-Request-ID", requestID)
		c.Next()
	}
}

// Logger logs request details
func (m *Middleware) Logger() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		// Read the request body
		var bodyBytes []byte
		if c.Request.Body != nil {
			bodyBytes, _ = io.ReadAll(c.Request.Body)
			c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		}

		c.Next()

		// Calculate request duration
		duration := time.Since(start)

		// Get request ID
		requestID, _ := c.Get("request_id")

		// Log request details
		m.logger.Info("request completed",
			logger.String("request_id", requestID.(string)),
			logger.String("method", c.Request.Method),
			logger.String("path", path),
			logger.String("query", raw),
			logger.Int("status", c.Writer.Status()),
			logger.String("duration", duration.String()),
			logger.String("ip", c.ClientIP()),
			logger.String("user_agent", c.Request.UserAgent()),
		)
	}
}

// Recovery handles panic recovery
func (m *Middleware) Recovery() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				requestID, _ := c.Get("request_id")
				m.logger.Error("panic recovered",
					logger.String("request_id", requestID.(string)),
					logger.String("error", err.(string)),
				)
				c.AbortWithStatusJSON(500, gin.H{
					"error": "Internal server error",
				})
			}
		}()
		c.Next()
	}
}
