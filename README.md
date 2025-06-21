# HostWatch

A command-line Go application that pings a given hostname or IP address (IPv4 or IPv6) until it replies to the ping.

## Features

- Supports both IPv4 and IPv6 addresses
- Accepts hostnames and resolves them automatically
- Continuously pings until the host responds
- Shows detailed ping information including round-trip time
- Graceful error handling and informative messages

## Usage

```bash
hostwatch.exe <hostname or IP address>
```

### Examples

```bash
# Ping a hostname
hostwatch.exe google.com

# Ping an IPv4 address
hostwatch.exe 8.8.8.8

# Ping an IPv6 address
hostwatch.exe 2001:4860:4860::8888
```

## Requirements

- Windows 11 (or compatible Windows version)
- Administrator privileges (required for raw ICMP socket operations)

## Important Notes

- This application requires elevated privileges to create raw ICMP sockets
- Run as Administrator or the application will fail with permission errors
- The application will continuously ping until it receives a response
- Press Ctrl+C to stop the application at any time

## How It Works

1. **Address Resolution**: The application first resolves the provided hostname or validates the IP address
2. **Protocol Detection**: Automatically detects whether to use IPv4 or IPv6 based on the resolved address
3. **Continuous Pinging**: Sends ICMP Echo Request packets every second until a reply is received
4. **Response Verification**: Validates that received replies match the sent requests
5. **Success Notification**: Reports when the host becomes responsive

## Output Example

```
Watching host: google.com
Pinging until host responds... (Press Ctrl+C to stop)
Resolved to: 142.250.191.14:0 (IPv4)
Ping 1 to 142.250.191.14:0: timeout or error (read ip4 0.0.0.0: i/o timeout)
Ping 2 to 142.250.191.14:0: timeout or error (read ip4 0.0.0.0: i/o timeout)
PING reply from 142.250.191.14: seq=3 time=23.456ms (IPv4)
Host google.com is now responding!
```

## Building from Source

If you want to build the application yourself:

```bash
go mod tidy
go build -o hostwatch.exe
```

## Improvements Over Original Code

- Added command-line argument parsing
- Added IPv6 support alongside IPv4
- Implemented continuous ping loop until success
- Enhanced error handling and user feedback
- Added hostname resolution with fallback between IPv4/IPv6
- Improved code organization and documentation
- Added proper type conversions for ICMP message parsing
