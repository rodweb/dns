package resolver

import (
	msg "github.com/rodweb/dns/internal/message"
)

type Resolver interface {
	Resolve(message *msg.Message) (*msg.Message, error)
}

func New(resolverAddress string) (Resolver, error) {
	if resolverAddress == "" {
		return NewDefaultResolver(), nil
	}

	return NewForwardingResolver(resolverAddress)
}
