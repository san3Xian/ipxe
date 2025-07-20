package httpd

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net/http"
	"time"
)

// Server provides simple file server with logging

type Server struct {
	root string
}

func New(root string) *Server { return &Server{root: root} }

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/status", s.status)
	fs := http.FileServer(http.Dir(s.root))
	mux.Handle("/", s.logMiddleware(fs))
	log.Infof("HTTP server listening on %s", addr)
	return http.ListenAndServe(addr, mux)
}

func (s *Server) logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		lrw := &logResponse{ResponseWriter: w}
		next.ServeHTTP(lrw, r)
		log.Infof("%s - - [%s] \"%s %s %s\" %d %d \"%s\" \"%s\"",
			r.RemoteAddr,
			start.Format("02/Jan/2006:15:04:05 -0700"),
			r.Method,
			r.URL.Path,
			r.Proto,
			lrw.status,
			lrw.size,
			r.Referer(),
			r.UserAgent(),
		)
	})
}

type logResponse struct {
	http.ResponseWriter
	status int
	size   int
}

func (l *logResponse) WriteHeader(code int) {
	l.status = code
	l.ResponseWriter.WriteHeader(code)
}

func (l *logResponse) Write(b []byte) (int, error) {
	if l.status == 0 {
		l.status = 200
	}
	n, err := l.ResponseWriter.Write(b)
	l.size += n
	return n, err
}

func (s *Server) status(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "OK")
}
