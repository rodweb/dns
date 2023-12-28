package message

import (
	"encoding/binary"
	"fmt"
)

// Header is a struct that represents a DNS internal header
// The header is 12 bytes long
// https://tools.ietf.org/html/rfc1035#section-4.1
type Header struct {
	// ID represents the Package identifier (ID).
	// A random identifier is assigned to query packets and
	// the response packages must reply with the same ID.
	ID uint16
	// IsResponse represents the Query/Response flag (QR).
	// 0 = Query, 1 = Response
	IsResponse bool
	// OperationCode represents the Operation code (OPCODE).
	// It is 4 bits long and typically is 0 (standard query).
	OperationCode OperationCode
	// AuthoritativeAnswer represents the Authoritative Answer (AA) flag.
	// If set to 1, the responding server is an authority for
	// the domain name in question section. This means it "owns" the domain.
	AuthoritativeAnswer bool
	// Truncated represents the Truncation flag (TC).
	// If set to 1, the internal was truncated because it exceeded 512 bytes.
	// In that case, the query must be repeated using TCP.
	Truncated bool
	// RecursionDesired represents the Recursion Desired (RD) flag.
	// If set to 1, the client wants the server to recursively resolve the query.
	RecursionDesired bool
	// RecursionAvailable represents the Recursion Available (RA) flag.
	// If set to 1, the server supports recursive queries.
	RecursionAvailable bool
	// Reserved represents the Reserved (Z) flag.
	// It is 3 bits long and was originally reserved for later use, but now
	// used for DNSSEC queries.
	Reserved uint8
	// ResponseCode represents the Response code (RCODE).
	// It is 4 bits long and is set by the server to indicate the status of the query.
	// 0 = No error condition
	ResponseCode ResponseCode
	// QuestionCount represents the number of entries in the question section (QDCOUNT).
	QuestionCount uint16
	// AnswerCount represents the number of entries in the answer section (ANCOUNT).
	AnswerCount uint16
	// AuthorityCount represents the number of entries in the authority records section (NSCOUNT).
	AuthorityCount uint16
	// AdditionalCount represents the number of entries in the additional records section (ARCOUNT).
	AdditionalCount uint16
}

type OperationCode uint8

const (
	Query OperationCode = iota
)

type ResponseCode uint8

const (
	Succeeded      ResponseCode = 0
	NotImplemented ResponseCode = 4
)

func (h *Header) Bytes() []byte {
	result := make([]byte, 12)

	// ID
	binary.BigEndian.PutUint16(result[0:2], h.ID)
	// Flags (IsResponse, OperationCode, Authoritative, Truncated, RD, RA, Z, RCODE)
	flags := uint16(0)
	if h.IsResponse {
		flags |= 1 << 15
	}
	flags |= uint16(h.OperationCode) << 11
	if h.AuthoritativeAnswer {
		flags |= 1 << 10
	}
	if h.Truncated {
		flags |= 1 << 9
	}
	if h.RecursionDesired {
		flags |= 1 << 8
	}
	if h.RecursionAvailable {
		flags |= 1 << 7
	}
	flags |= uint16(h.Reserved) << 4
	flags |= uint16(h.ResponseCode)

	binary.BigEndian.PutUint16(result[2:4], flags)
	binary.BigEndian.PutUint16(result[4:6], h.QuestionCount)
	binary.BigEndian.PutUint16(result[6:8], h.AnswerCount)
	binary.BigEndian.PutUint16(result[8:10], h.AuthorityCount)
	binary.BigEndian.PutUint16(result[10:12], h.AdditionalCount)

	return result
}

// String returns a string representation of the Header struct
func (h *Header) String() string {
	return fmt.Sprintf("Header{ID: %d, OPCODE: %d, QDCOUNT: %d, ANCOUNT: %d}",
		h.ID,
		h.OperationCode,
		h.QuestionCount,
		h.AnswerCount,
	)
}

// headerFromBytes decodes the DNS message header from the message packet
func headerFromBytes(packet []byte, offset *int) *Header {
	flags := binary.BigEndian.Uint16(packet[2:4])
	header := &Header{
		ID:                  binary.BigEndian.Uint16(packet[0:2]),
		IsResponse:          (flags >> 15 & 0x01) != 0,
		OperationCode:       OperationCode((flags >> 11)) & 0x0F,
		AuthoritativeAnswer: (flags >> 10 & 0x01) != 0,
		Truncated:           (flags >> 9 & 0x01) != 0,
		RecursionDesired:    (flags >> 8 & 0x01) != 0,
		RecursionAvailable:  (flags >> 7 & 0x01) != 0,
		Reserved:            uint8((flags >> 4)) & 0x07,
		QuestionCount:       binary.BigEndian.Uint16(packet[4:6]),
		AnswerCount:         binary.BigEndian.Uint16(packet[6:8]),
		AuthorityCount:      binary.BigEndian.Uint16(packet[8:10]),
		AdditionalCount:     binary.BigEndian.Uint16(packet[10:12]),
	}
	*offset += 12
	fmt.Println(header.String())
	return header
}

func GetResponseCode(header *Header) ResponseCode {
	// Standard query
	if header.OperationCode == Query {
		return Succeeded
	}
	return NotImplemented
}
