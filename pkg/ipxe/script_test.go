package ipxe

import (
	"dpxe/pkg/config"
	"testing"
)

func TestScript(t *testing.T) {
	c := &config.Config{Boot: config.Sanboot, ISOPath: "/a.iso"}
	s := Script(c)
	if s == "" {
		t.Fatal("script empty")
	}
}
