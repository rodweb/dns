package main

import (
	"flag"
	"fmt"
	msg "github.com/rodweb/dns/internal/message"
	rsv "github.com/rodweb/dns/internal/resolver"
	"net"
	"strconv"
	"strings"
)

func main() {
	resolverPtr := flag.String("resolver", "", "Resolver forward requests to (address:port)")
	flag.Parse()
	var resolver rsv.Resolver
	fmt.Printf("Resolver: %s\n", *resolverPtr)
	if *resolverPtr == "" {
		resolver = rsv.NewDefaultResolver()
	} else {
		parts := strings.Split(*resolverPtr, ":")
		if len(parts) != 2 {
			fmt.Println("Invalid resolver address")
			return
		}
		ip := net.ParseIP(parts[0])
		if ip == nil {
			fmt.Println("Invalid resolver address")
			return
		}
		port, err := strconv.Atoi(parts[1])
		if err != nil {
			fmt.Println("Invalid resolver address")
			return
		}
		resolver = &rsv.ForwardingResolver{
			IP:   ip,
			Port: port,
		}
	}

	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		fmt.Println("Failed to resolve UDP address:", err)
		return
	}

	udpConn, err := net.ListenUDP("udp", udpAddr)
	if err != nil {
		fmt.Println("Failed to bind to address:", err)
		return
	}
	defer udpConn.Close()

	buf := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buf)
		if err != nil {
			fmt.Println("Error receiving data:", err)
			break
		}

		var s strings.Builder
		s.WriteString("Packet received:\n")
		for i, b := range buf[:size] {
			if (i % 8) == 0 {
				s.WriteString("\n")
			}
			s.WriteString(fmt.Sprintf("0x%02x, ", b))
		}
		s.WriteString("\n\n")
		fmt.Printf(s.String())
		replyBytes, err := HandleReply(buf[:size], resolver)
		if err != nil {
			fmt.Println("Failed to handle reply:", err)
			continue
		}

		_, err = udpConn.WriteToUDP(replyBytes, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func HandleReply(data []byte, resolver rsv.Resolver) ([]byte, error) {
	message := &msg.Message{}
	err := message.Decode(data)
	if err != nil {
		fmt.Println("Failed to decode internal:", err)
		return nil, err
	}
	reply, err := resolver.Resolve(message)
	if err != nil {
		fmt.Println("Failed to resolve internal:", err)
		return nil, err
	}
	return reply.Encode(), nil
}
