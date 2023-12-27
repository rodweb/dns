package resolver

import (
	"fmt"
	cfg "github.com/rodweb/dns/internal/config"
	msg "github.com/rodweb/dns/internal/message"
	"strconv"
	"strings"
)

type DefaultResolver struct {
	dnsRecords map[string]*cfg.Record
}

func NewDefaultResolver(dnsRecords map[string]*cfg.Record) *DefaultResolver {
	return &DefaultResolver{
		dnsRecords: dnsRecords,
	}
}

func (r *DefaultResolver) Resolve(request *msg.Message) (*msg.Message, error) {
	response := newResponse(request)
	for _, question := range request.Questions {
		recordType := toRecordType(question.Type)
		// TODO: forward queries with non implemented but valid record types
		if recordType == "" {
			return nil, fmt.Errorf("invalid record type %d", question.Type)
		}

		// TODO: forward queries when not found locally
		key := fmt.Sprintf("%s:%s", recordType, question.Name)
		record, ok := r.dnsRecords[key]
		if !ok {
			continue
		}

		answer, err := newAnswer(question, record)
		if err != nil {
			// TODO: respond with a valid DNS message
			return nil, err
		}

		response.Questions = append(response.Questions, question)
		response.Answers = append(response.Answers, answer)
	}

	response.Header.QuestionCount = uint16(len(response.Questions))
	response.Header.AnswerCount = uint16(len(response.Answers))
	// TODO: handle unanswered questions
	response.Header.ResponseCode = msg.GetResponseCode(request.Header)

	return response, nil
}

func newResponse(req *msg.Message) *msg.Message {
	return &msg.Message{
		Header: &msg.Header{
			ID:               req.Header.ID,
			IsResponse:       true,
			OperationCode:    req.Header.OperationCode,
			RecursionDesired: req.Header.RecursionDesired,
		},
		Questions: make([]*msg.Question, 0, req.Header.QuestionCount),
		Answers:   make([]*msg.Answer, 0, req.Header.QuestionCount),
	}
}

func newAnswer(q *msg.Question, record *cfg.Record) (*msg.Answer, error) {
	recordType := toRecordType(q.Type)
	data, err := parseRecordData(recordType, record.Value)
	if err != nil {
		return nil, err
	}

	// TODO: why not copy Type and Class from Question?
	return &msg.Answer{
		Name:  q.Name,
		Type:  1,
		Class: 1,
		TTL:   uint32(record.TTL),
		Data:  data,
	}, nil
}

// TODO: move this to the config.Record itself?
func parseRecordData(recordType string, data string) ([]byte, error) {
	switch recordType {
	case "A":
		return parseARecord(data)
	default:
		return nil, fmt.Errorf("failed to parse record data for type %s", recordType)
	}
}

func parseARecord(data string) ([]byte, error) {
	ipParts := strings.Split(data, ".")
	if len(ipParts) != 4 {
		return nil, fmt.Errorf("invalid IP data")
	}
	ip := make([]byte, 4)
	for i, part := range ipParts {
		value, err := strconv.Atoi(part)
		if err != nil {
			return nil, fmt.Errorf("invalid IP data")
		}
		ip[i] = byte(value)
	}
	return ip, nil
}

func toRecordType(questionType uint16) string {
	if questionType == 1 {
		return "A"
	}
	return ""
}
