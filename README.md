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
# ISO image that will be exposed to iPXE clients
iso_path: /path/to/iso

# Directory served over HTTP
http_root: /var/www

# iPXE script returned at /boot.ipxe
ipxe_script: /path/to/boot.ipxe

dhcp:
  # Beginning of lease range
  start_ip: 192.168.0.100
  # Last address in the range
  end_ip: 192.168.0.200
  # Default gateway handed to clients
  router: 192.168.0.1
  # DNS server for clients
  dns: 8.8.8.8
  # Lease time in seconds
  lease_duration: 60
```

## Example: booting Proxmox VE

If you have the `proxmox-ve_8.4-1.iso` file in `/isos`, you can expose it with
the following configuration:

```yaml
iso_path: /isos/proxmox-ve_8.4-1.iso
http_root: /isos
ipxe_script: /etc/dpxe/proxmox.ipxe
dhcp:
  start_ip: 10.0.0.100
  end_ip: 10.0.0.200
  router: 10.0.0.1
  dns: 8.8.8.8
  lease_duration: 300
```

The accompanying `proxmox.ipxe` script might look like:

```text
#!ipxe
initrd http://${serverip}/proxmox-ve_8.4-1.iso
chain memdisk iso raw
```

## Tests

Run `go test ./...` to execute all tests.
