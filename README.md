# dpxe

dpxe is a small iPXE and DHCP server written in Go. It exposes a simple HTTP file server and replies to DHCPv4 requests. Configuration is provided via YAML.

## Usage

```
# start server with configuration
./dpxe -config dpxe.yaml

# query current leases
./dpxe -client
```

## Configuration

```yaml
iso_path: /path/to/iso
http_root: /var/www
ipxe_script: /path/to/boot.ipxe
dhcp:
  start_ip: 192.168.0.100
  end_ip: 192.168.0.200
  router: 192.168.0.1
  dns: 8.8.8.8
  lease_duration: 60
```

## Tests

Run `go test ./...` to execute all tests.
