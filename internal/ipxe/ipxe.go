package ipxe

import (
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
)

// Handler serves a static iPXE script.
type Handler struct {
	scriptPath string
}

// New returns a handler that serves the given iPXE script file.
func New(script string) *Handler {
	return &Handler{scriptPath: script}
}

// ServeHTTP writes the iPXE script to the client.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadFile(h.scriptPath)
	if err != nil {
		log.Error(err)
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
	w.Header().Set("Content-Type", "text/plain")
	w.Write(data)
}
