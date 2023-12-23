package main

import (
	"encoding/binary"
	"fmt"
	"net"
	// Uncomment this block to pass the first stage
	// "net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.
	fmt.Println("Logs from your program will appear here!")

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

		receivedData := string(buf[:size])
		fmt.Printf("Received %d bytes from %s: %s\n", size, source, receivedData)

		// Create an empty response
		message := Message{
			Header: Header{
				ID:      1234,
				QR:      true,
				OPCODE:  0,
				AA:      false,
				TC:      false,
				RD:      false,
				RA:      false,
				Z:       0,
				RCODE:   0,
				QDCOUNT: 0,
				ANCOUNT: 0,
				NSCOUNT: 0,
				ARCOUNT: 0,
			},
		}
		response := message.Serialize()
		fmt.Printf("Sending %d bytes to %s: %s\n", len(response), source, response)

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func (m *Message) Serialize() []byte {
	message := make([]byte, 12)

	// ID
	binary.BigEndian.PutUint16(message[0:2], m.Header.ID)

	// Flags (QR, OPCODE, AA, TC, RD, RA, Z, RCODE)
	flags := uint16(0)

	if m.Header.QR {
		flags |= 1 << 15
	}

	flags |= uint16(m.Header.OPCODE) << 11

	if m.Header.AA {
		flags |= 1 << 10
	}

	if m.Header.TC {
		flags |= 1 << 9
	}

	if m.Header.RD {
		flags |= 1 << 8
	}

	if m.Header.RA {
		flags |= 1 << 7
	}

	flags |= uint16(m.Header.Z) << 4
	flags |= uint16(m.Header.RCODE)

	binary.BigEndian.PutUint16(message[2:4], flags)

	return message
}

// Message is a struct that represents a DNS message
type Message struct {
	Header Header
}

// Header is a struct that represents a DNS message header
// The header is 12 bytes long
type Header struct {
	ID      uint16 // Package identifier
	QR      bool   // Query/Response flag
	OPCODE  uint8  // Operation code - 4bits
	AA      bool   // Authoritative Answer
	TC      bool   // Truncation flag
	RD      bool   // Recursion Desired
	RA      bool   // Recursion Available
	Z       uint8  // Reserved for future use - 3bits
	RCODE   uint8  // Response code - 4bits
	QDCOUNT uint16 // Number of entries in the question section
	ANCOUNT uint16 // Number of resource records in the answer section
	NSCOUNT uint16 // Number of name server resource records in the authority records section
	ARCOUNT uint16 // Number of resource records in the additional records section
}
