package requestlog

import (
	"errors"
	"io"
	"net"
	"net/http"
	"time"
)

type logEntry struct {
	RequestID         string
	ReceivedTime      time.Time
	RequestMethod     string
	RequestURL        string
	RequestHeaderSize int64
	RequestBodySize   int64
	UserAgent         string
	Referer           string
	Proto             string

	RemoteIP string
	ServerIP string

	Status             int
	ResponseHeaderSize int64
	ResponseBodySize   int64
	Latency            time.Duration
}

func ipFromHostPort(hp string) string {
	h, _, err := net.SplitHostPort(hp)
	if err != nil {
		return ""
	}

	if len(h) > 0 && h[0] == '[' {
		return h[1 : len(h)-1]
	}

	return h
}

type readCounterCloser struct {
	r   io.ReadCloser
	n   int64
	err error
}

func (rcc *readCounterCloser) Read(p []byte) (n int, err error) {
	if rcc.err != nil {
		return 0, rcc.err
	}
	n, rcc.err = rcc.r.Read(p)
	rcc.n += int64(n)
	return n, rcc.err
}

func (rcc *readCounterCloser) Close() error {
	rcc.err = errors.New("read from closed reader")
	return rcc.r.Close()
}

type writeCounter int64

func (wc *writeCounter) Write(p []byte) (n int, err error) {
	*wc += writeCounter(len(p))
	return len(p), nil
}

func headerSize(h http.Header) int64 {
	var wc writeCounter
	e := h.Write(&wc)
	if e != nil {
		println(e)
	}
	return int64(wc) + 2 // CRLF
}

type responseStats struct {
	w     http.ResponseWriter
	hsize int64
	wc    writeCounter
	code  int
}

func (r *responseStats) Header() http.Header {
	return r.w.Header()
}

func (r *responseStats) WriteHeader(statusCode int) {
	if r.code != 0 {
		return
	}

	r.hsize = headerSize(r.w.Header())
	r.w.WriteHeader(statusCode)
	r.code = statusCode
}

func (r *responseStats) Write(p []byte) (n int, err error) {
	if r.code == 0 {
		r.WriteHeader(http.StatusOK)
	}

	n, err = r.w.Write(p)
	_, e := r.wc.Write(p)
	if e != nil {
		println(e)
	}
	return
}

func (r *responseStats) size() (hdr, body int64) {
	if r.code == 0 {
		return headerSize(r.w.Header()), 0
	}
	return r.hsize, int64(r.wc)
}
