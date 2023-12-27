package main

import (
	"log"
	"net"
)

// Listener listens for DNS requests
type Listener struct {
	handler *Handler
	udpAddr *net.UDPAddr
}

// NewListener creates a new Listener
func NewListener(handler *Handler) *Listener {
	udpAddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:2053")
	if err != nil {
		log.Fatalln("Failed to resolve UDP address:", err)
	}

	return &Listener{
		handler: handler,
		udpAddr: udpAddr,
	}
}

// ListenAndServe starts the Listener
func (l *Listener) ListenAndServe() error {
	udpConn, err := net.ListenUDP("udp", l.udpAddr)
	if err != nil {
		log.Fatalln("Failed to bind to address:", err)
	}
	defer func() {
		err := udpConn.Close()
		if err != nil {
			log.Println("Failed to close UDP connection:", err)
		}
	}()

	buffer := make([]byte, 512)

	for {
		size, source, err := udpConn.ReadFromUDP(buffer)
		if err != nil {
			log.Println("Error receiving data:", err)
			continue
		}

		response, err := l.handler.Handle(buffer[:size])
		if err != nil {
			log.Println("Failed to handle packet:", err)
			continue
		}

		_, err = udpConn.WriteToUDP(response, source)
		if err != nil {
			log.Println("Failed to send response:", err)
			continue
		}
	}
}
