package message

import (
	"bytes"
	"encoding/binary"
	"fmt"
)

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

func (q *Question) Bytes() []byte {
	var buff bytes.Buffer

	buff.Write(serializeDomainName(q.Name))
	binary.Write(&buff, binary.BigEndian, q.Type)
	binary.Write(&buff, binary.BigEndian, q.Class)

	return buff.Bytes()
}

// String returns a string representation of the Question struct
func (q *Question) String() string {
	return fmt.Sprintf("Question{Name: %s, Type: %d, Class: %d}",
		q.Name,
		q.Type,
		q.Class,
	)
}

// questionsFromBytes decodes the DNS message question section from the message packet
func questionsFromBytes(data []byte, offset *int, count uint16) []*Question {
	result := make([]*Question, count)
	for i := 0; i < int(count); i++ {
		question := questionFromBytes(data, offset)
		result[i] = question
	}
	return result
}

// questionFromBytes decodes a DNS question from question section of the message packet
func questionFromBytes(data []byte, offset *int) *Question {
	name := domainNameFromBytes(data, offset)
	question := &Question{
		Name:  name,
		Type:  binary.BigEndian.Uint16(data[*offset : *offset+2]),
		Class: binary.BigEndian.Uint16(data[*offset+2 : *offset+4]),
	}
	*offset += 4
	fmt.Println(question.String())
	return question
}
