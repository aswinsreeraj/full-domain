package middleware

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type ctxKey string

const RequestIDKey ctxKey = "requestID"

func RequestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := uuid.NewString()

		c.Set(string(RequestIDKey), id)

		ctx := context.WithValue(c.Request.Context(), RequestIDKey, id)
		c.Request = c.Request.WithContext(ctx)

		c.Writer.Header().Set("X-Request-ID", id)

		c.Next()
	}
}
