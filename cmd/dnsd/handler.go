package main

import (
	"fmt"
	"github.com/rodweb/dns/internal/config"
	msg "github.com/rodweb/dns/internal/message"
	rsv "github.com/rodweb/dns/internal/resolver"
	"log"
	"strings"
)

// Handler is a DNS query handler.
type Handler struct {
	dnsRecords map[string]string
	resolver   rsv.Resolver
}

// NewHandler creates a new Handler.
func NewHandler(records []config.Record) *Handler {
	dnsRecords := make(map[string]string)

	// Map DNS records to be served by the DNS server
	for _, r := range records {
		key := fmt.Sprintf("%s:%s", r.Type, r.Name)
		dnsRecords[key] = r.Value
	}

	resolver := rsv.NewDefaultResolver()

	return &Handler{
		dnsRecords: dnsRecords,
		resolver:   resolver,
	}
}

// Handle handles a DNS query.
func (h *Handler) Handle(packet []byte) ([]byte, error) {
	printPacket(packet)

	// Parse the DNS message
	message, err := parseMessage(packet)
	if err != nil {
		log.Println("Failed to parse message:", err)
		return nil, err
	}

	// Resolve the DNS queries
	response, err := h.resolver.Resolve(message)
	if err != nil {
		log.Println("Failed to resolve:", err)
		return nil, err
	}

	return response.Bytes(), nil
}

// printPacket pretty prints the UDP packet
func printPacket(packet []byte) {
	var s strings.Builder
	s.WriteString("Packet received:")
	for i, b := range packet {
		if (i % 8) == 0 {
			s.WriteString("\n")
		}
		s.WriteString(fmt.Sprintf("0x%02x, ", b))
	}
	s.WriteString("\n")
	log.Println(s.String())
}

// parseMessage parses a DNS message.
func parseMessage(packet []byte) (*msg.Message, error) {
	// TODO: Handle errors without panicking
	return msg.MessageFromBytes(packet), nil
}
