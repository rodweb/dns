package message

import (
	"bytes"
	"encoding/binary"
	"strings"
)

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

func domainNameFromBytes(data []byte, offset *int) string {
	var result []string
	for {
		// If first two bits are 1, it's a pointer
		if ((data[*offset] >> 6) & 0x3) == 0x3 {
			nameOffset := int((binary.BigEndian.Uint16(data[*offset:*offset+2]) << 2) >> 2)
			result = append(result, domainNameFromBytes(data, &nameOffset))
			*offset += 2
			break
		}

		length := int(data[*offset])
		*offset++
		label := string(data[*offset : *offset+length])
		result = append(result, label)
		*offset += length

		if data[*offset] == 0x00 {
			*offset += 1
			break
		}
	}
	return strings.Join(result, ".")
}
