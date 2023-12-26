package resolver

import (
	msg "github.com/rodweb/dns/internal/message"
)

type DefaultResolver struct{}

func NewDefaultResolver() *DefaultResolver {
	return &DefaultResolver{}
}

func (r *DefaultResolver) Resolve(request *msg.Message) (reply *msg.Message, error error) {
	reply = newReply(request)
	for i, question := range request.Questions {
		reply.Questions[i] = &msg.Question{
			Name:  question.Name,
			Type:  1,
			Class: 1,
		}
		reply.Answers[i] = &msg.Answer{
			Name:  question.Name,
			Type:  1,
			Class: 1,
			TTL:   60,
			Data:  []byte{0x8, 0x8, 0x8, 0x8},
		}
	}
	return reply, nil
}

func newReply(req *msg.Message) *msg.Message {
	return &msg.Message{
		Header: &msg.Header{
			ID:                  req.Header.ID,
			IsResponse:          true,
			OperationCode:       req.Header.OperationCode,
			AuthoritativeAnswer: false,
			Truncated:           false,
			RecursionDesired:    req.Header.RecursionDesired,
			RecursionAvailable:  false,
			Reserved:            0,
			ResponseCode:        msg.GetResponseCode(req.Header),
			QueryCount:          req.Header.QueryCount,
			AnswerCount:         req.Header.QueryCount,
			AuthorityCount:      0,
			AdditionalCount:     0,
		},
		Questions: make([]*msg.Question, req.Header.QueryCount),
		Answers:   make([]*msg.Answer, req.Header.QueryCount),
	}
}
