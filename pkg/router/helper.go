package router

import (
	"context"
	"fmt"
)

func PathValue(ctx context.Context, name string) string {
	value := ctx.Value(ctxKey(name))
	if value == nil {
		return fmt.Errorf("no value found for %s in context", name).Error()
	}
	return value.(string)
}
