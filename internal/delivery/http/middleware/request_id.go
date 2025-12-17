package middleware

import (
	"full-domain/pkg/woodpecker"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()
		c.Writer.Header().Set("X-Request-ID", id)

		reqLogger := woodpecker.NewRequestLogger(id)

		ctx := woodpecker.WithLogger(c.Request.Context(), reqLogger)
		c.Request = c.Request.WithContext(ctx)

		c.Next()
	}
}
