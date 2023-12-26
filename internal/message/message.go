package message

import (
	"bytes"
	"fmt"
)

// Message is a struct that represents a DNS internal
type Message struct {
	Header    *Header
	Questions []*Question
	Answers   []*Answer
}

func (m *Message) Encode() []byte {
	fmt.Printf("Encoding Header %+v\n", m.Header)
	headerBytes := m.Header.Encode()
	var resourceRecordBytes []byte
	for i, question := range m.Questions {
		fmt.Printf("Encoding Question %d: %+v\n", i, question)
		resourceRecordBytes = append(resourceRecordBytes, question.Encode()...)
	}
	for i, answer := range m.Answers {
		fmt.Printf("Encoding Answer %d: %+v\n", i, answer)
		resourceRecordBytes = append(resourceRecordBytes, answer.Encode()...)
	}
	return bytes.Join([][]byte{headerBytes, resourceRecordBytes}, []byte{})
}

func (m *Message) Decode(data []byte) error {
	m.Header = headerFromBytes(data[:12])
	offset := 12
	fmt.Printf("Decoding Header %+v\n", m.Header)
	m.Questions = questionsFromBytes(data, &offset, m.Header.QueryCount)
	m.Answers = answersFromBytes(data, &offset, m.Header.AnswerCount)
	return nil
}
