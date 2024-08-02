package middleware

import (
	"net/http"

	ctxUtil "employee-worklog-service/utils/ctx"
	"github.com/rs/xid"
)

const requestIdHeaderKey = "X-Request-ID"

func RequestId(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		requestId := r.Header.Get(requestIdHeaderKey)
		if requestId == "" {
			requestId = xid.New().String()
		}

		ctx = ctxUtil.SetRequestId(ctx, requestId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
