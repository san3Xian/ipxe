package httpd

import (
	"io/ioutil"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/example/dpxe/internal/ipxe"
)

func TestServer(t *testing.T) {
	tmpDir, err := ioutil.TempDir("", "webroot")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tmpDir)
	ioutil.WriteFile(tmpDir+"/index.html", []byte("ok"), 0644)

	ipxeFile, _ := ioutil.TempFile("", "script.ipxe")
	ipxeFile.WriteString("#!ipxe")
	ipxeFile.Close()

	srv := New("127.0.0.1:8888", tmpDir, ipxe.New(ipxeFile.Name()))
	go srv.Serve()
	time.Sleep(100 * time.Millisecond)

	resp, err := http.Get("http://" + srv.Addr() + "/index.html")
	if err != nil {
		t.Fatal(err)
	}
	b, _ := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if string(b) != "ok" {
		t.Fatalf("unexpected %s", b)
	}
}
