// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

package middleware

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/amankapoor/placementpal/internal/platform/web"
)

type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *loggingResponseWriter {
	return &loggingResponseWriter{w, http.StatusOK}
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

// ReqResLogger writes some information about the request to the logs in
// the format: TraceID : (200) GET /foo -> IP ADDR (latency)
func ReqResLogger(next web.Handler) web.Handler {

	// Wrap this handler around the next one provided.
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, params web.Params) error {
		v := ctx.Value(web.KeyValues).(*web.Values)
		//next(ctx, w, r, params)

		lrw := NewLoggingResponseWriter(w)
		next(ctx, lrw, r, params)

		statusCode := lrw.statusCode

		log.Printf("(%d) : %s %s -> %s (%s) <-- %d : %s : %s",
			v.StatusCode,
			r.Method, r.URL.Path,
			r.RemoteAddr, time.Since(v.Now),
			statusCode,
			v.TraceID,
			r.UserAgent(),
		)

		// This is the top of the food chain. At this point all error
		// handling has been done including logging.
		return nil
	}
}
