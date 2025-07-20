package ipxe

import (
	"io/ioutil"
	"net/http"
	"path/filepath"

	log "github.com/sirupsen/logrus"
)

// Handler returns a http.Handler that serves the iPXE script stored in rootDir/pxeFile
func Handler(rootDir, pxeFile string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		path := filepath.Join(rootDir, pxeFile)
		data, err := ioutil.ReadFile(path)
		if err != nil {
			log.Errorf("failed to read ipxe file: %v", err)
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/plain")
		w.Write(data)
	})
}
