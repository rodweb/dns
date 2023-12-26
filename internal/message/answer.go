package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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
	// Length of the Data field in bytes (RDLENGTH)
	Length uint16
	// Data specific to the query type (RDATA)
	Data []byte
}

func (a Answer) Encode() []byte {
	var buff bytes.Buffer

	buff.Write(serializeDomainName(a.Name))
	binary.Write(&buff, binary.BigEndian, a.Type)
	binary.Write(&buff, binary.BigEndian, a.Class)
	binary.Write(&buff, binary.BigEndian, a.TTL)
	binary.Write(&buff, binary.BigEndian, uint16(len(a.Data)))
	buff.Write(a.Data)

	return buff.Bytes()[:buff.Len()]
}

func answersFromBytes(data []byte, offset *int, count uint16) []*Answer {
	result := make([]*Answer, count)
	for i := 0; i < int(count); i++ {
		answer := answerFromBytes(data, offset)
		fmt.Printf("Decoding Answer %d: %+v\n", i, answer)
		result[i] = answer
	}
	return result
}

func answerFromBytes(data []byte, offset *int) *Answer {
	name := domainNameFromBytes(data, offset)
	answer := &Answer{
		Name:   name,
		Type:   binary.BigEndian.Uint16(data[*offset : *offset+2]),
		Class:  binary.BigEndian.Uint16(data[*offset+2 : *offset+4]),
		TTL:    binary.BigEndian.Uint32(data[*offset+4 : *offset+8]),
		Length: binary.BigEndian.Uint16(data[*offset+8 : *offset+10]),
	}
	*offset += 10
	// TODO: Handle different record types
	answer.Data = data[*offset : *offset+int(answer.Length)]
	*offset += int(answer.Length)
	return answer
}
