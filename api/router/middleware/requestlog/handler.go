package requestlog

import (
	"io"
	"net"
	"net/http"
	"time"

	ctxUtil "employee-worklog-service/utils/ctx"
	"github.com/rs/zerolog"
)

type Handler struct {
	handler http.Handler
	logger  *zerolog.Logger
}

func NewHandler(h http.HandlerFunc, l *zerolog.Logger) *Handler {
	return &Handler{
		handler: h,
		logger:  l,
	}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	start := time.Now()

	le := &logEntry{
		RequestID:         ctxUtil.RequestId(r.Context()),
		ReceivedTime:      start,
		RequestMethod:     r.Method,
		RequestURL:        r.URL.String(),
		RequestHeaderSize: headerSize(r.Header),
		UserAgent:         r.UserAgent(),
		Referer:           r.Referer(),
		Proto:             r.Proto,
		RemoteIP:          ipFromHostPort(r.RemoteAddr),
	}

	if addr, ok := r.Context().Value(http.LocalAddrContextKey).(net.Addr); ok {
		le.ServerIP = ipFromHostPort(addr.String())
	}

	r2 := new(http.Request)
	*r2 = *r
	rcc := &readCounterCloser{r: r.Body}
	r2.Body = rcc
	w2 := &responseStats{w: w}

	h.handler.ServeHTTP(w2, r2)

	le.Latency = time.Since(start)
	if rcc.err == nil && rcc.r != nil {
		// If the handler hasn't encountered an error in the body (like EOF),
		// then consume the rest of the body to provide an accurate rcc.n.
		_, e := io.Copy(io.Discard, rcc)
		if e != nil {
			println("Reached EOF or error")
		}
	}
	le.RequestBodySize = rcc.n
	le.Status = w2.code
	if le.Status == 0 {
		le.Status = http.StatusOK
	}
	le.ResponseHeaderSize, le.ResponseBodySize = w2.size()
	h.logger.Info().
		Str("request_id", le.RequestID).
		Time("recieved_time", le.ReceivedTime).
		Str("method", le.RequestMethod).
		Str("url", le.RequestURL).
		Int64("header_size", le.RequestHeaderSize).
		Int64("body_size", le.RequestBodySize).
		Str("agent", le.UserAgent).
		Str("referer", le.Referer).
		Str("proto", le.Proto).
		Str("remote_ip", le.RemoteIP).
		Str("server_ip", le.ServerIP).
		Int("status", le.Status).
		Int64("resp_header_size", le.ResponseHeaderSize).
		Int64("resp_body_size", le.ResponseBodySize).
		Dur("latency", le.Latency).
		Msg("")
}
