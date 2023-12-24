package main

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"net"
	"strings"
	// Uncomment this block to pass the first stage
	// "net"
)

func main() {
	// You can use print statements as follows for debugging, they'll be visible when running tests.

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

		message := &Message{}
		err = message.Decode(buf[:size])
		if err != nil {
			fmt.Println("Failed to decode message:", err)
			continue
		}

		reply := newReply(message)
		for i, question := range message.Questions {
			reply.Questions[i] = &Question{
				Name:  question.Name,
				Type:  1,
				Class: 1,
			}
			reply.Answers[i] = &Answer{
				Name:  question.Name,
				Type:  1,
				Class: 1,
				TTL:   60,
				RDATA: []byte{0x8, 0x8, 0x8, 0x8},
			}
		}
		replyBytes := reply.Encode()

		_, err = udpConn.WriteToUDP(replyBytes, source)
		if err != nil {
			fmt.Println("Failed to send response:", err)
		}
	}
}

func (m *Message) Decode(data []byte) error {
	m.Header = headerFromBytes(data[:12])
	fmt.Printf("Decoding Header %+v\n", m.Header)
	m.Questions = questionsFromBytes(data[12:], m.Header.QueryCount)
	fmt.Printf("Decoding Question %+v\n", m.Questions[0])
	return nil
}

func headerFromBytes(data []byte) *Header {
	flags := binary.BigEndian.Uint16(data[2:4])
	return &Header{
		ID:                  binary.BigEndian.Uint16(data[0:2]),
		IsResponse:          (flags >> 15 & 0x01) != 0,
		OperationCode:       uint8((flags >> 11)) & 0x0F,
		AuthoritativeAnswer: (flags >> 10 & 0x01) != 0,
		Truncated:           (flags >> 9 & 0x01) != 0,
		RecursionDesired:    (flags >> 8 & 0x01) != 0,
		RecursionAvailable:  (flags >> 7 & 0x01) != 0,
		Reserved:            uint8((flags >> 4)) & 0x07,
		QueryCount:          binary.BigEndian.Uint16(data[4:6]),
		AnswerCount:         binary.BigEndian.Uint16(data[6:8]),
		AuthoritativeCount:  binary.BigEndian.Uint16(data[8:10]),
		AdditionalCount:     binary.BigEndian.Uint16(data[10:12]),
	}
}

func questionsFromBytes(data []byte, count uint16) []*Question {
	result := make([]*Question, count)
	var offset int
	for i := 0; i < int(count); i++ {
		question, bytesRead := questionFromBytes(data[offset:])
		if bytesRead == 0 {
			break
		}
		result[i] = question
		offset += bytesRead
	}
	return result
}

func questionFromBytes(data []byte) (*Question, int) {
	name, offset := domainNameFromBytes(data)
	return &Question{
		Name:  name,
		Type:  binary.BigEndian.Uint16(data[offset : offset+2]),
		Class: binary.BigEndian.Uint16(data[offset+2 : offset+4]),
	}, offset + 4
}

func domainNameFromBytes(data []byte) (string, int) {
	var result []string
	var offset int
	for {
		bytesToRead := int(data[offset])
		offset++
		label := string(data[offset : offset+bytesToRead])
		result = append(result, label)
		offset += bytesToRead
		if data[offset] == 0x00 {
			break
		}
	}
	return strings.Join(result, "."), int(offset)
}

func newReply(req *Message) *Message {
	return &Message{
		Header: &Header{
			ID:                  req.Header.ID,
			IsResponse:          true,
			OperationCode:       req.Header.OperationCode,
			AuthoritativeAnswer: false,
			Truncated:           false,
			RecursionDesired:    req.Header.RecursionDesired,
			RecursionAvailable:  false,
			Reserved:            0,
			ResponseCode:        getResponseCode(req.Header),
			QueryCount:          req.Header.QueryCount,
			AnswerCount:         req.Header.QueryCount,
			AuthoritativeCount:  0,
			AdditionalCount:     0,
		},
		Questions: make([]*Question, req.Header.QueryCount),
		Answers:   make([]*Answer, req.Header.QueryCount),
	}
}

func getResponseCode(header *Header) uint8 {
	// Standard query (opcode == 0)
	if header.OperationCode == 0 {
		return 0
	}
	// Not implemented
	return 4
}

// Message is a struct that represents a DNS message
type Message struct {
	Header    *Header
	Questions []*Question
	Answers   []*Answer
}

func (m *Message) Encode() []byte {
	fmt.Printf("Encoding Header %+v\n", m.Header)
	fmt.Printf("Encoding Question %+v\n", m.Questions[0])
	headerBytes := m.Header.Serialize()
	var resourceRecordBytes []byte
	for _, question := range m.Questions {
		resourceRecordBytes = append(resourceRecordBytes, question.Encode()...)
	}
	for _, answer := range m.Answers {
		resourceRecordBytes = append(resourceRecordBytes, answer.Encode()...)
	}
	return bytes.Join([][]byte{headerBytes, resourceRecordBytes}, []byte{})
}

// Header is a struct that represents a DNS message header
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
	OperationCode uint8
	// AuthoritativeAnswer represents the Authoritative Answer (AA) flag.
	// If set to 1, the responding server is an authority for
	// the domain name in question section. This means it "owns" the domain.
	AuthoritativeAnswer bool
	// Truncated represents the Truncation flag (TC).
	// If set to 1, the message was truncated because it exceeded 512 bytes.
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
	ResponseCode uint8
	// QueryCount represents the number of entries in the question section (QDCOUNT).
	QueryCount uint16
	// AnswerCount represents the number of entries in the answer section (ANCOUNT).
	AnswerCount uint16
	// AuthoritativeCount represents the number of entries in the authority records section (NSCOUNT).
	AuthoritativeCount uint16
	// AdditionalCount represents the number of entries in the additional records section (ARCOUNT).
	AdditionalCount uint16
}

func (h *Header) Serialize() []byte {
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
	binary.BigEndian.PutUint16(result[4:6], h.QueryCount)
	binary.BigEndian.PutUint16(result[6:8], h.AnswerCount)
	binary.BigEndian.PutUint16(result[8:10], h.AuthoritativeCount)
	binary.BigEndian.PutUint16(result[10:12], h.AdditionalCount)

	return result
}

// Question is a struct that represents a DNS question
// https://www.rfc-editor.org/rfc/rfc1035#section-4.1.2
type Question struct {
	// Domain name represented as a sequence of labels
	// which are encoded as <length><label> where <length> is a single octet
	Name string
	// Type of the query (1 = A, 2 = NS, 5 = CNAME, 6 = SOA, 12 = PTR, 15 = MX, 16 = TXT)
	// https://www.rfc-editor.org/rfc/rfc1035#section-3.2.2
	Type uint16
	// Class of the query (1 = IN, 2 = CS, 3 = CH, 4 = HS)
	// https://www.rfc-editor.org/rfc/rfc1035#section-3.2.4
	Class uint16
}

func (q *Question) Encode() []byte {
	var buff bytes.Buffer

	buff.Write(serializeDomainName(q.Name))
	binary.Write(&buff, binary.BigEndian, q.Type)
	binary.Write(&buff, binary.BigEndian, q.Class)

	return buff.Bytes()
}

// Answer is a struct that represents a DNS answer
// https://www.rfc-editor.org/rfc/rfc1035#section-3.2.1
type Answer struct {
	Name string
	// Type of the query (1 = A, 2 = NS, 5 = CNAME, 6 = SOA, 12 = PTR, 15 = MX, 16 = TXT)
	// https://www.rfc-editor.org/rfc/rfc1035#section-3.2.2
	Type uint16
	// Class of the query (1 = IN, 2 = CS, 3 = CH, 4 = HS)
	// https://www.rfc-editor.org/rfc/rfc1035#section-3.2.4
	Class uint16
	// Time to live in seconds
	// The duration that the RR can be cached before querying the DNS server again
	TTL uint32
	// Length of the RDATA field in bytes
	RDLENTH uint16
	// Data specific to the query type
	RDATA []byte
}

func (a Answer) Encode() []byte {
	var buff bytes.Buffer

	buff.Write(serializeDomainName(a.Name))
	binary.Write(&buff, binary.BigEndian, a.Type)
	binary.Write(&buff, binary.BigEndian, a.Class)
	binary.Write(&buff, binary.BigEndian, a.TTL)
	binary.Write(&buff, binary.BigEndian, uint16(len(a.RDATA)))
	buff.Write(a.RDATA)

	return buff.Bytes()[:buff.Len()]
}

func serializeDomainName(domain string) []byte {
	var buff bytes.Buffer

	labels := strings.Split(domain, ".")
	for _, label := range labels {
		buff.WriteByte(byte(len(label)))
		buff.Write([]byte(label))
	}
	buff.WriteByte(0x00)

	return buff.Bytes()
}
