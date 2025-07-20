# dpxe

`dpxe` is a simple Go based DHCP and iPXE server. It reads ISO boot
information from a YAML configuration file and exposes a small API via a
unix socket to query current DHCP leases.

## Building

```
go build ./cmd/dpxe
```

## Running

Start the server:

```
./dpxe -config config.yaml
```

Query lease information from another process:

```
./dpxe -client
```

The iPXE script can be downloaded from `http://<server>:8080/boot.ipxe`.
The DHCP server listens on UDP port 6767.

## Configuration

Example `config.yaml`:

```yaml
isos:
  - name: demo
    kernel: http://example.com/vmlinuz
    initrd: http://example.com/initrd.img
    cmdline: console=ttyS0
```
