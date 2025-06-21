package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

func main() {
	// Check command line arguments
	if len(os.Args) < 2 {
		fmt.Printf("Usage: %s <hostname or IP address>\n\n", os.Args[0])
		fmt.Println("Examples:")
		fmt.Println("  hostwatch.exe google.com")
		fmt.Println("  hostwatch.exe 8.8.8.8")
		fmt.Println("  hostwatch.exe 2001:4860:4860::8888")
		os.Exit(1)
	}

	host := os.Args[1]

	fmt.Printf("Watching host: %s\n", host)

	// Determine if this is IPv4 or IPv6
	ip := net.ParseIP(host)
	var isIPv6 bool
	var dst net.Addr
	var err error
	var isHostname = false

	if ip != nil {
		// It's already an IP address
		isIPv6 = ip.To4() == nil
		if isIPv6 {
			dst, err = net.ResolveIPAddr("ip6", host)
		} else {
			dst, err = net.ResolveIPAddr("ip4", host)
		}
	} else {
		// It's a hostname, try to resolve it
		isHostname = true

		// First try IPv4
		dst, err = net.ResolveIPAddr("ip4", host)
		if err != nil {
			// If IPv4 fails, try IPv6
			dst, err = net.ResolveIPAddr("ip6", host)
			if err == nil {
				isIPv6 = true
			}
		}
	}

	if err != nil {
		fmt.Printf("Error resolving address '%s': %v\n", host, err)
		os.Exit(1)
	}

	if isHostname {
        // The user provided a hostname instead of an IP address, so display the resolved address.
		fmt.Printf("Resolved to: %s (%s)\n", dst, getIPVersion(isIPv6))
	}

	// Keep pinging until we get a successful response
	fmt.Println("\nPinging until host responds... (Press Ctrl+C to stop)")

	seq := 1
	for {
		success := ping(dst, seq, isIPv6)
		if success {
			fmt.Printf("\nHost %s is now responding!\n", host)
			break
		}

		seq++
		time.Sleep(1 * time.Second) // Wait 1 second between pings
	}
}

func getIPVersion(isIPv6 bool) string {
	if isIPv6 {
		return "IPv6"
	}
	return "IPv4"
}

func ping(dst net.Addr, seq int, isIPv6 bool) bool {
	var conn *icmp.PacketConn
	var err error
	var icmpType icmp.Type
	var protocol string

	if isIPv6 {
		conn, err = icmp.ListenPacket("ip6:ipv6-icmp", "::")
		icmpType = ipv6.ICMPTypeEchoRequest
		protocol = "IPv6"
	} else {
		conn, err = icmp.ListenPacket("ip4:icmp", "0.0.0.0")
		icmpType = ipv4.ICMPTypeEcho
		protocol = "IPv4"
	}

	if err != nil {
		fmt.Printf("Error creating ICMP connection (%s): %v\n", protocol, err)
		return false
	}
	defer conn.Close()

	// Create an ICMP Echo Request message
	message := &icmp.Message{
		Type: icmpType,
		Code: 0,
		Body: &icmp.Echo{
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: []byte("hostwatch ping"),
		},
	}

	// Marshal the ICMP message into bytes
	var messageBytes []byte
	if isIPv6 {
		messageBytes, err = message.Marshal(nil)
	} else {
		messageBytes, err = message.Marshal(nil)
	}

	if err != nil {
		fmt.Printf("Error marshaling ICMP message: %v\n", err)
		return false
	}

	// Record the time we're sending the packet
	start := time.Now()

	// Send the ICMP packet to the destination
	_, err = conn.WriteTo(messageBytes, dst)
	if err != nil {
		fmt.Printf("Error sending ICMP packet to %v: %v\n", dst, err)
		return false
	}

	// Set a read timeout
	err = conn.SetReadDeadline(time.Now().Add(3 * time.Second))
	if err != nil {
		fmt.Printf("Error setting read deadline: %v\n", err)
		return false
	}

	// Buffer to receive the reply packet
	reply := make([]byte, 1500)

	// Read the ICMP reply from the network
	n, peer, err := conn.ReadFrom(reply)
	if err != nil {
		fmt.Printf("Ping %d to %v: timeout or error (%v)\n", seq, dst, err)
		return false
	}

	// Calculate the round-trip time
	duration := time.Since(start)

	// Parse the received ICMP message using the correct function signature
	// Protocol numbers: 1 for ICMPv4, 58 for ICMPv6
	var replyMessage *icmp.Message
	if isIPv6 {
		replyMessage, err = icmp.ParseMessage(58, reply[:n]) // ICMPv6 protocol number
	} else {
		replyMessage, err = icmp.ParseMessage(1, reply[:n]) // ICMPv4 protocol number
	}

	if err != nil {
		// If parsing fails, it's likely corrupted or unknown ICMP message
		fmt.Printf("Ping %d to %v: failed to parse ICMP message from %v: %v\n",
			seq, dst, peer, err)
		return false
	}

	// Check if this is an Echo Reply message
	var expectedEchoReplyType icmp.Type
	if isIPv6 {
		expectedEchoReplyType = ipv6.ICMPTypeEchoReply
	} else {
		expectedEchoReplyType = ipv4.ICMPTypeEchoReply
	}

	if replyMessage.Type != expectedEchoReplyType {
		// This is some other ICMP message (like destination unreachable, time exceeded, etc.)
		// These are error messages, not successful ping responses
		fmt.Printf("Ping %d to %v: received ICMP %v from %v\n",
			seq, dst, replyMessage.Type, peer)
		return false
	}

	// It's an Echo Reply - verify it matches our request
	if echoReply, ok := replyMessage.Body.(*icmp.Echo); ok {
		expectedID := os.Getpid() & 0xffff
		if echoReply.ID == expectedID && echoReply.Seq == seq {
			fmt.Printf("PING reply from %v: seq=%d time=%v (%s)\n",
				peer, seq, duration, protocol)
			return true
		} else {
			fmt.Printf("Ping %d to %v: received Echo Reply with wrong ID/seq (ID=%d, seq=%d) from %v\n",
				seq, dst, echoReply.ID, echoReply.Seq, peer)
			return false
		}
	}

	// If we get here, it's not a proper Echo Reply format
	fmt.Printf("Ping %d to %v: received malformed Echo Reply from %v\n",
		seq, dst, peer)
	return false
}
