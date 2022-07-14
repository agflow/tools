package cache

import (
	"context"
	"time"
)

// Service declares the interface of a cache service
type Service interface {
	Get(context.Context, string) ([]byte, error)
	Set(context.Context, string, interface{}, time.Duration) error
}
