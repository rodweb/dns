package resolver

import (
	"fmt"
	msg "github.com/rodweb/dns/internal/message"
	"math/rand"
	"net"
	"sync"
)

type ForwardingResolver struct {
	IP   net.IP
	Port int
}

// TODO: Improve error handling
func (r *ForwardingResolver) Resolve(originalRequest *msg.Message) (reply *msg.Message, error error) {
	// When forwarding a originalRequest, we need to split the questions into multiple requests

	// We map an ID to a question
	questionMap := make(map[uint16]*msg.Question)

	var wg sync.WaitGroup

	resChan := make(chan *msg.Message, len(originalRequest.Questions))

	for _, question := range originalRequest.Questions {
		id := generateID()
		questionMap[id] = question
		req := &msg.Message{
			Header: &msg.Header{
				ID:            id,
				OperationCode: originalRequest.Header.OperationCode,
				QueryCount:    1,
			},
			Questions: []*msg.Question{
				question,
			},
		}
		wg.Add(1)

		go func() {
			res, err := makeRequest(r.IP, r.Port, req.Encode())
			if err != nil {
				fmt.Println("Failed to make originalRequest:", err)
				return
			}
			resMessage := &msg.Message{}
			err = resMessage.Decode(res)
			if err != nil {
				fmt.Println("Failed to decode response:", err)
				return
			}
			resChan <- resMessage
			wg.Done()
		}()
	}

	fmt.Printf("IdMap: %+v\n", questionMap)

	go func() {
		wg.Wait()
		close(resChan)
	}()

	questions := make([]*msg.Question, 0, originalRequest.Header.QueryCount)
	answers := make([]*msg.Answer, 0, originalRequest.Header.QueryCount)
	for res := range resChan {
		if res.Header.ResponseCode != 0 {
			fmt.Println("Response code is not 0")
			continue
		}
		question, ok := questionMap[res.Header.ID]
		if !ok {
			fmt.Println("ID not found in map")
			continue
		}
		if len(res.Answers) == 0 {
			fmt.Println("No answers")
			continue
		}
		questions = append(questions, question)
		answers = append(answers, res.Answers[0])
	}

	reply = &msg.Message{
		Header: &msg.Header{
			ID:               originalRequest.Header.ID,
			IsResponse:       true,
			RecursionDesired: originalRequest.Header.RecursionDesired,
			OperationCode:    originalRequest.Header.OperationCode,
			ResponseCode:     msg.GetResponseCode(originalRequest.Header),
			QueryCount:       uint16(len(questions)),
			AnswerCount:      uint16(len(answers)),
		},
		Questions: questions,
		Answers:   answers,
	}

	return reply, nil
}

// generateID generates a random number between 0 and 65535
func generateID() uint16 {
	return uint16(rand.Intn(65535))
}

func makeRequest(ip net.IP, port int, data []byte) ([]byte, error) {
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
