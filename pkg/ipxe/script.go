package ipxe

import (
	"dpxe/pkg/config"
	"fmt"
)

// Script returns ipxe script content based on config
func Script(c *config.Config) string {
	switch c.Boot {
	case config.Sanboot:
		return fmt.Sprintf("#!ipxe\ninitrd %s\nsanboot \\${initrd}\n", c.ISOPath)
	case config.Memdisk:
		return fmt.Sprintf("#!ipxe\nkernel %s\ninitrd %s\nboot\n", c.Kernel, c.Initrd)
	default:
		return "#!ipxe\n"
	}
}
