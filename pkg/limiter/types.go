package limiter

import "golang.org/x/net/context"

type Limiter interface {
	// Limit 是否出发限流。 true: 触发限流
	Limit(ctx context.Context, key string) (bool, error)
}
