package ipxe

import (
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

func TestServeHTTP(t *testing.T) {
	tmp, err := ioutil.TempFile("", "script*.ipxe")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmp.Name())
	tmp.WriteString("#!ipxe")
	tmp.Close()

	h := New(tmp.Name())
	req := httptest.NewRequest(http.MethodGet, "/boot.ipxe", nil)
	rr := httptest.NewRecorder()
	h.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Fatalf("status %d", rr.Code)
	}
	if rr.Body.String() != "#!ipxe" {
		t.Fatalf("body %s", rr.Body.String())
	}
}
