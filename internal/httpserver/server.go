package httpserver

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// New returns configured HTTP server serving from mux and logging requests
func New(addr string, handler http.Handler) *http.Server {
	return &http.Server{
		Addr:    addr,
		Handler: loggingMiddleware(handler),
	}
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		lrw := &logResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		start := time.Now()
		next.ServeHTTP(lrw, r)
		duration := time.Since(start)
		log.Infof("%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\" %v",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.RequestURI,
			r.Proto,
			lrw.statusCode,
			lrw.written,
			r.Referer(),
			r.UserAgent(),
			duration,
		)
	})
}

type logResponseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int
}

func (lrw *logResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *logResponseWriter) Write(b []byte) (int, error) {
	n, err := lrw.ResponseWriter.Write(b)
	lrw.written += n
	return n, err
}
