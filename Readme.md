# dpxe

`dpxe` is a small DHCP and iPXE helper written in Go. It provides a minimal IPv4 DHCP server, an embedded HTTP server for iPXE boot resources and a simple unix socket interface for querying current leases.

## Building

```
GOOS=darwin GOARCH=arm64 go build ./cmd/dpxe
```

## Running

Create a `config.yaml` file. Example:

```yaml
http_root: ./resources
pxe_file: boot.ipxe
socket: /tmp/dpxe.sock
dhcp:
  interface: eth0
  start_ip: 192.168.1.100
  end_ip: 192.168.1.110
  router: 192.168.1.1
  dns: 8.8.8.8
  lease_time: 1h
  subnet_mask: 255.255.255.0
```

Place any ISO images or additional files in the directory pointed by `http_root`. The `pxe_file` should be located relative to this directory and is served when a client boots via iPXE.

Run the server:

```
./dpxe -config config.yaml
```

To query current leases while the server is running:

```
./dpxe -config config.yaml -client
```

### Proxmox-VE 8.4 example

With Proxmox installed, configure a VM to boot from network (PXE). Provide an iPXE script similar to:

```
#!ipxe
kernel http://<dpxe_ip>:8080/vmlinuz
initrd http://<dpxe_ip>:8080/initrd.img
boot
```

Place kernel and initrd inside `http_root` directory.

## Testing

```
GOOS=darwin GOARCH=arm64 go test ./...
```

