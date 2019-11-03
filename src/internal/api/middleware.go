package api

import (
	"net"
	"net/http"
	"strings"
	"time"

	"github.com/howtocards/preview-queue/src/internal/config"
	"github.com/powerman/structlog"
)

type middlewareFunc func(http.Handler) http.Handler

// Provide a logger configured using request's context.
//
// Usually it should be first middleware.
func makeLogger(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			log := structlog.New(
				config.LogRemote, r.RemoteAddr,
				config.LogHTTPStatus, "",
				config.LogHTTPMethod, r.Method,
				config.LogFunc, strings.TrimPrefix(r.URL.Path, basePath),
			)
			r = r.WithContext(structlog.NewContext(r.Context(), log))

			next.ServeHTTP(w, r)
		})
	}
}

// go-swagger responders panic on error while writing response to client,
// this shouldn't result in crash - unlike a real, reasonable panic.
//
// Usually it should be second middleware (after logger).
func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			const code = http.StatusInternalServerError
			switch err := recover(); err := err.(type) {
			default:
				log := structlog.FromContext(r.Context(), nil)
				log.PrintErr("panic", config.LogHTTPStatus, code, "err", err, structlog.KeyStack, structlog.Auto)
				w.WriteHeader(code)
			case nil:
			case net.Error:
				log := structlog.FromContext(r.Context(), nil)
				log.PrintErr("recovered", config.LogHTTPStatus, code, "err", err)
				w.WriteHeader(code)
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func makeAccessLog(basePath string) middlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := wrapResponseWriter(w)

			next.ServeHTTP(ww, r)

			code := ww.StatusCode()

			log := structlog.FromContext(r.Context(), nil)
			if code < 500 {
				log.Info("handled", "in", time.Since(start), config.LogHTTPStatus, code)
			} else {
				log.PrintErr("failed to handle", config.LogHTTPStatus, code)
			}
		})
	}
}
