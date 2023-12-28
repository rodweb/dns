package message

import (
	"bytes"
)

// Message is a struct that represents a DNS internal
type Message struct {
	Header    *Header
	Questions []*Question
	Answers   []*Answer
}

// Bytes returns a byte array representation of the DNS message
func (m *Message) Bytes() []byte {
	var buffer bytes.Buffer
	buffer.Write(m.Header.Bytes())
	for _, question := range m.Questions {
		buffer.Write(question.Bytes())
	}
	for _, answer := range m.Answers {
		buffer.Write(answer.Bytes())
	}
	return buffer.Bytes()
}

// FromBytes decodes a DNS message from a byte array
func FromBytes(packet []byte) *Message {
	var offset int
	message := &Message{}
	message.Header = headerFromBytes(packet, &offset)
	message.Questions = questionsFromBytes(packet, &offset, message.Header.QuestionCount)
	message.Answers = answersFromBytes(packet, &offset, message.Header.AnswerCount)
	return message
}
