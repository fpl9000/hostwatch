# HostWatch

A command-line Go application that pings a given hostname or IP address (IPv4 or IPv6) until it replies to the ping.

## Features

- Works without Administrator privileges.
- Supports both IPv4 and IPv6 addresses.
- Accepts hostnames and resolves them automatically.
- Continuously pings until the host responds.
- Shows detailed ping information including round-trip time.
- Graceful error handling and informative messages.

## Usage

```
hostwatch { HOSTNAME | IPV4ADDRESS | IPV6ADDRESS }
```

### Examples

```
# Ping a hostname
hostwatch.exe google.com

# Ping an IPv4 address
hostwatch.exe 8.8.8.8

# Ping an IPv6 address
hostwatch.exe 2001:4860:4860::8888
```

## Output Example

```
$ ./hostwatch www.google.com
Watching host: www.google.com
Resolved to: 142.251.40.196 (IPv4)

Pinging until host responds... (Press Ctrl+C to stop)
PING reply from 142.251.40.196: seq=1 time=227.1973ms (IPv4)

Host www.google.com is now responding!
```

## Requirements

- Windows 11 (or compatible Windows version)

## Important Notes

- The application will continuously ping until it receives a response.
- Press Ctrl+C to stop the application before the host responds.

## Building from Source

If you want to build the application yourself:

```bash
go mod tidy
go build -o hostwatch.exe
```
