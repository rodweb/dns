package resolver

import (
	"fmt"
	msg "github.com/rodweb/dns/internal/message"
	"math/rand"
	"net"
	"strconv"
	"strings"
	"sync"
)

// ForwardingResolver is a resolver that forwards requests to another resolver
type ForwardingResolver struct {
	IP   net.IP
	Port int
}

// NewForwardingResolver creates a new forwarding resolver
func NewForwardingResolver(resolverAddress string) (*ForwardingResolver, error) {
	parts := strings.Split(resolverAddress, ":")
	if len(parts) != 2 {
		return nil, fmt.Errorf("invalid resolver address")
	}
	ip := net.ParseIP(parts[0])
	if ip == nil {
		return nil, fmt.Errorf("invalid IP address")
	}
	port, err := strconv.Atoi(parts[1])
	if err != nil {
		return nil, fmt.Errorf("invalid port")
	}
	return &ForwardingResolver{
		IP:   ip,
		Port: port,
	}, nil
}

// TODO: Improve error handling
// Resolve resolves a request by forwarding it to another resolver
func (r *ForwardingResolver) Resolve(originalMessage *msg.Message) (*msg.Message, error) {
	// When forwarding a message, we need to split the questions into multiple queries

	// Create a map of ID to question
	questionMap := make(map[uint16]*msg.Question)

	var wg sync.WaitGroup

	responseChan := make(chan *msg.Message, len(originalMessage.Questions))

	// For each question, create a new query and forward it to the resolver
	for _, question := range originalMessage.Questions {
		id := generateID()
		questionMap[id] = question
		query := &msg.Message{
			Header: &msg.Header{
				ID:            id,
				OperationCode: originalMessage.Header.OperationCode,
				QuestionCount: 1,
			},
			Questions: []*msg.Question{
				question,
			},
		}
		wg.Add(1)

		go func(name string) {
			defer wg.Done()
			fmt.Printf("Forwarding query for %s\n", name)
			response, err := forwardQuery(r.IP, r.Port, query.Bytes())
			if err != nil {
				fmt.Println("Failed forward query:", err)
				return
			}
			responseChan <- msg.FromBytes(response)
		}(question.Name)
	}

	// Wait for all responses to be received
	go func() {
		wg.Wait()
		close(responseChan)
	}()

	questions := make([]*msg.Question, 0, originalMessage.Header.QuestionCount)
	answers := make([]*msg.Answer, 0, originalMessage.Header.QuestionCount)

	// For each response, add the question and answer to the original response
	for response := range responseChan {
		if response.Header.ResponseCode != 0 {
			fmt.Println("Response code is not 0")
			continue
		}
		question, ok := questionMap[response.Header.ID]
		if !ok {
			fmt.Println("ID not found in map")
			continue
		}
		if len(response.Answers) == 0 {
			fmt.Println("No answers")
			continue
		}
		questions = append(questions, question)
		answers = append(answers, response.Answers[0])
	}

	return &msg.Message{
		Header: &msg.Header{
			ID:               originalMessage.Header.ID,
			IsResponse:       true,
			RecursionDesired: originalMessage.Header.RecursionDesired,
			OperationCode:    originalMessage.Header.OperationCode,
			ResponseCode:     msg.GetResponseCode(originalMessage.Header),
			QuestionCount:    uint16(len(questions)),
			AnswerCount:      uint16(len(answers)),
		},
		Questions: questions,
		Answers:   answers,
	}, nil
}

// generateID generates a random number between 0 and 65535
func generateID() uint16 {
	return uint16(rand.Intn(65535))
}

// TODO: Reuse UDP connections
// forwardQuery forwards a query to another resolver
func forwardQuery(ip net.IP, port int, data []byte) ([]byte, error) {
	addr, err := net.ResolveUDPAddr("udp", fmt.Sprintf("%s:%d", ip.String(), port))
	if err != nil {
		return nil, err
	}

	conn, err := net.DialUDP("udp", nil, addr)
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	_, err = conn.Write(data)

	buffer := make([]byte, 512)
	size, _, err := conn.ReadFromUDP(buffer)
	if err != nil {
		return nil, err
	}

	return buffer[:size], nil
}
