package clickhouse

import (
	"context"
	"net/http"

	"gorm.io/gorm"
)

type CustomContext struct {
	Database *gorm.DB
}

var customContextKey = "CUSTOM_CONTEXT"

func CreateContext(args *CustomContext, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		customContext := &CustomContext{
			Database: args.Database,
		}
		requestWithCtx := r.WithContext(context.WithValue(r.Context(), customContextKey, customContext))
		next.ServeHTTP(w, requestWithCtx)
	})
}

func GetContext(ctx context.Context) *CustomContext {
	customContext, ok := ctx.Value(customContextKey).(*CustomContext)
	if !ok {
		return nil
	}
	return customContext
}
