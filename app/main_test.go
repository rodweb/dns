package main

import (
	msg "github.com/rodweb/dns/internal/message"
	"testing"
)

func TestDecodeMessage(t *testing.T) {
	packet := []byte{
		0x00, 0x01, // ID
		0x01, 0x20, // Flags
		0x00, 0x01, // Question count
		0x00, 0x00, // Answer count
		0x00, 0x00, // Authority count
		0x00, 0x00, // Additional count
		0x06, 0x67, 0x6F, 0x6F, 0x67, 0x6C, 0x65, // google
		0x03, 0x63, 0x6F, 0x6D, // com
		0x00,       // End of domain name
		0x00, 0x01, // Type
		0x00, 0x01, // Class
	}
	message := &msg.Message{}
	err := message.Decode(packet)
	if err != nil {
		t.Error("Failed to decode internal:", err)
	}
	if message.Header.ID != 1 {
		t.Error("Failed to decode ID")
	}
	if message.Header.IsResponse != false {
		t.Error("Failed to decode QR")
	}
	if message.Header.OperationCode != 0 {
		t.Error("Failed to decode OPCODE")
	}
	if message.Header.AuthoritativeAnswer != false {
		t.Error("Failed to decode AA")
	}
	if message.Header.Truncated != false {
		t.Error("Failed to decode TC")
	}
	if message.Header.RecursionDesired != true {
		t.Error("Failed to decode RD")
	}
	if message.Header.RecursionAvailable != false {
		t.Error("Failed to decode RA")
	}
	if message.Header.Reserved != 2 {
		t.Error("Failed to decode Z")
	}
	if message.Header.ResponseCode != 0 {
		t.Error("Failed to decode RCODE")
	}
	if message.Header.QueryCount != 1 {
		t.Error("Failed to decode QDCOUNT")
	}
	if message.Header.AnswerCount != 0 {
		t.Error("Failed to decode ANCOUNT")
	}
	if message.Header.AuthorityCount != 0 {
		t.Error("Failed to decode NSCOUNT")
	}
	if message.Header.AdditionalCount != 0 {
		t.Error("Failed to decode ARCOUNT")
	}
}

func TestDecodeCompressedMessage(t *testing.T) {
	packet := []byte{
		0x00, 0x01, // ID
		0x01, 0x00, // Flags
		0x00, 0x02, // Question count
		0x00, 0x00, // Answer count
		0x00, 0x00, // Authority count
		0x00, 0x00, // Additional count
		0x03, 0x61, 0x62, 0x63, // abc
		0x11, 0x6c, 0x6f, 0x6e, 0x67, 0x61, 0x73, 0x73, 0x64, 0x6f, 0x6d, 0x61, 0x69, 0x6e, 0x6e, 0x61, 0x6d, 0x65, // longassdomainname
		0x03, 0x63, 0x6f, 0x6d, // com
		0x00,       // End of domain name
		0x00, 0x01, // Type
		0x00, 0x01, // Class
		0x03, 0x64, 0x65, 0x66, // def
		0xc0, 0x10, // Pointer to longassdomainname
		0x00, 0x01, // Type
		0x00, 0x01, // Class
	}
	message := &msg.Message{}
	err := message.Decode(packet)
	if err != nil {
		t.Error("Failed to decode internal:", err)
	}
	if len(message.Questions) != 2 {
		t.Error("Failed to decode questions")
	}
	if message.Questions[0].Name != "abc.longassdomainname.com" {
		t.Error("Failed to decode first question name")
	}
	if message.Questions[1].Name != "def.longassdomainname.com" {
		t.Error("Failed to decode second question name")
	}
}
