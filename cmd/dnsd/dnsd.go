package main

import (
	"github.com/rodweb/dns/internal/config"
	"log"
)

func main() {
	err := config.Load()
	if err != nil {
		log.Fatalln("Failed to load config:", err)
	}

	handler := NewHandler(config.Get().Records)
	listener := NewListener(handler)

	err = listener.ListenAndServe()
	if err != nil {
		log.Fatalln("Failed to start listener:", err)
	}
}
