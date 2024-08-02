package middleware

import "net/http"

const (
	HeaderKeyContentType   = "Content-Type"
	HeaderValueContentType = "application/json;charset=utf-8"
)

func ContentTypeJSON(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set(HeaderKeyContentType, HeaderValueContentType)
		next.ServeHTTP(w, r)
	})
}
