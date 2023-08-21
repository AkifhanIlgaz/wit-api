package ctx

import "context"

type key string

const (
	uidKey key = "uid"
)

func WithUid(ctx context.Context, uid string) context.Context {
	return context.WithValue(ctx, uidKey, uid)
}

func Uid(ctx context.Context) *string {
	val := ctx.Value(uidKey)
	if uid, ok := val.(string); ok {
		return &uid
	}
	return nil
}
