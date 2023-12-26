package resolver

import msg "github.com/rodweb/dns/internal/message"

type Resolver interface {
	Resolve(message *msg.Message) (*msg.Message, error)
}
