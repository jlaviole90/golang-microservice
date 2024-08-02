package ctx

import "context"

const keyRequestId key = "requestId"

type key string

func RequestId(ctx context.Context) string {
	requestId, _ := ctx.Value(keyRequestId).(string)

	return requestId
}

func SetRequestId(ctx context.Context, requestId string) context.Context {
	return context.WithValue(ctx, keyRequestId, requestId)
}
