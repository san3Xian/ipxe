# dpxe

`dpxe` provides a minimal DHCP and HTTP server with ipxe script generation. It is
intended for offline installations of systems such as CentOS, Debian or Proxmox.

## Build

```
go build ./cmd/dpxe
```

For macOS arm64:

```
GOOS=darwin GOARCH=arm64 go build ./cmd/dpxe
```

Run tests:

```
go test ./...
GOOS=darwin GOARCH=arm64 go test ./...
```

## Configuration

Configuration is provided via `config.yaml`.

```
addr: ":8080"      # address for HTTP server
http_root: "./data" # directory to serve files
iso_path: "proxmox.iso" # path to ISO used in sanboot
kernel: "vmlinuz"   # kernel path for memdisk
initrd: "initrd.img" # initrd path for memdisk
boot: sanboot       # boot method: sanboot or memdisk
dhcp_range: 192.168.1.100-192.168.1.200
```

### Example: Proxmox VE 8.4

1. Place `proxmox-ve_8.4.iso` in `./data` directory.
2. Set `iso_path: proxmox-ve_8.4.iso` and `boot: sanboot` in `config.yaml`.
3. Start the server:

```
./dpxe -config config.yaml
```

PXE clients should boot using the generated `boot.ipxe` script.

## Querying status

While `dpxe` is running you can query status information:

```
./dpxe -client
```
