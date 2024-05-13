package utils

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"

	"github.com/lakhansamani/cloud-container/internal/middleware"
)

// GinContext to get gin context from context
func GinContextFromContext(ctx context.Context) (*gin.Context, error) {
	ginContext := ctx.Value(middleware.GinContextKeyValue)
	if ginContext == nil {
		err := fmt.Errorf("could not retrieve gin.Context")
		return nil, err
	}
	gc, ok := ginContext.(*gin.Context)
	if !ok {
		err := fmt.Errorf("gin.Context has wrong type")
		return nil, err
	}
	return gc, nil
}
