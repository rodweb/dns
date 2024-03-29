package main

import (
	"fmt"
	cfg "github.com/rodweb/dns/internal/config"
	msg "github.com/rodweb/dns/internal/message"
	rsv "github.com/rodweb/dns/internal/resolver"
	"log"
	"strings"
)

type Resolver interface {
	Resolve(request *msg.Message) (*msg.Message, error)
}

// Handler is a DNS query handler.
type Handler struct {
	dnsRecords map[string]*cfg.Record
	resolver   Resolver
}

// NewHandler creates a new Handler.
func NewHandler(records []*cfg.Record) *Handler {
	dnsRecords := make(map[string]*cfg.Record)

	// Map DNS records to be served by the DNS server
	for _, r := range records {
		key := fmt.Sprintf("%s:%s", r.Type, r.Name)
		dnsRecords[key] = r
	}

	resolver := rsv.NewDefaultResolver(dnsRecords)

	log.Printf("Resolver initialized with %d DNS records\n", len(dnsRecords))

	return &Handler{
		dnsRecords: dnsRecords,
		resolver:   resolver,
	}
}

// Handle handles a DNS query.
func (h *Handler) Handle(packet []byte) ([]byte, error) {
	printPacket(packet)

	// Parse the DNS request
	request, err := parseMessage(packet)
	if err != nil {
		log.Println("Failed to parse request:", err)
		return nil, err
	}

	// Resolve the DNS queries
	response, err := h.resolver.Resolve(request)
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
	return msg.FromBytes(packet), nil
}
