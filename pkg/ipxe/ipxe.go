package ipxe

import (
	"fmt"
	"net/http"

	"github.com/example/dpxe/pkg/config"
	log "github.com/sirupsen/logrus"
)

// Server represents the ipxe server
type Server struct {
	cfg *config.Config
}

// New creates a new ipxe server
func New(cfg *config.Config) *Server {
	return &Server{cfg: cfg}
}

// Serve starts the ipxe HTTP server
func (s *Server) Serve(addr string) error {
	http.HandleFunc("/boot.ipxe", s.bootScript)
	log.Infof("iPXE server listening on %s", addr)
	return http.ListenAndServe(addr, nil)
}

func (s *Server) bootScript(w http.ResponseWriter, r *http.Request) {
	if len(s.cfg.ISOs) == 0 {
		http.Error(w, "no isos defined", http.StatusInternalServerError)
		return
	}
	iso := s.cfg.ISOs[0]
	script := fmt.Sprintf("#!ipxe\nkernel %s %s\ninitrd %s\nboot\n", iso.Kernel, iso.Cmdline, iso.Initrd)
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte(script))
}
