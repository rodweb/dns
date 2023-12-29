package api

import (
	"log"
	"net/http"
)

func main() {
	registerRoutes()

	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatalln("failed to start api server", err)
	}
}
