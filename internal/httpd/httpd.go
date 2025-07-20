package httpd

import (
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"
)

// Server wraps an HTTP file server and optional iPXE handler.
type Server struct {
	addr        string
	root        string
	ipxeHandler http.Handler
}

// New creates a new HTTP server bound to addr serving files from root.
func New(addr, root string, ipxe http.Handler) *Server {
	return &Server{addr: addr, root: root, ipxeHandler: ipxe}
}

// Addr returns the configured listen address.
func (s *Server) Addr() string { return s.addr }

// Serve starts the HTTP server and blocks.
func (s *Server) Serve() error {
	mux := http.NewServeMux()
	if s.ipxeHandler != nil {
		mux.Handle("/boot.ipxe", s.ipxeHandler)
	}
	fileServer := http.FileServer(http.Dir(s.root))
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
		log.Infof("%s - - [%s] \"%s %s %s\" %d %d \"%s\"",
			r.RemoteAddr,
			time.Now().Format("02/Jan/2006:15:04:05 -0700"),
			r.Method, r.RequestURI, r.Proto,
			http.StatusOK,
			0,
			r.UserAgent(),
		)
	})
	log.Infof("HTTP server listening on %s serving %s", s.addr, s.root)
	return http.ListenAndServe(s.addr, mux)
}
