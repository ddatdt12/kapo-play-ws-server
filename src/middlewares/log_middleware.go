package middlewares

import (
	"bufio"
	"bytes"
	"errors"
	"log"
	"net"
	"net/http"
	"time"

	"fmt"

	"github.com/fatih/color"
	"github.com/gorilla/mux"
)

type LogResponseWriter struct {
	http.ResponseWriter
	statusCode int
	buf        bytes.Buffer
}

func NewLogResponseWriter(w http.ResponseWriter) *LogResponseWriter {
	return &LogResponseWriter{ResponseWriter: w}
}

func (w *LogResponseWriter) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

func (w *LogResponseWriter) Write(body []byte) (int, error) {
	w.buf.Write(body)
	return w.ResponseWriter.Write(body)
}

type LogMiddleware struct {
	logger *log.Logger
}

func NewLogMiddleware(logger *log.Logger) *LogMiddleware {
	return &LogMiddleware{logger: logger}
}

func (m *LogMiddleware) Func() mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			startTime := time.Now()

			logRespWriter := NewLogResponseWriter(w)
			next.ServeHTTP(logRespWriter, r)

			duration := time.Since(startTime)
			statusCode := logRespWriter.statusCode
			body := logRespWriter.buf.String()
			path := r.URL.Path

			// Colorize the log output
			durationColor := color.New(color.FgCyan).SprintFunc()
			statusCodeColor := color.New(color.FgYellow).SprintFunc()
			bodyColor := color.New(color.FgGreen).SprintFunc()
			pathColor := color.New(color.FgBlue).SprintFunc()

			m.logger.Printf(
				"duration=%s status=%s path=%s body=%s",
				durationColor(duration.String()),
				statusCodeColor(fmt.Sprintf("%d", statusCode)),
				pathColor(path),
				bodyColor(body),
			)
		})
	}
}

func (w *LogResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	h, ok := w.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, errors.New("hijack not supported")
	}
	return h.Hijack()
}
