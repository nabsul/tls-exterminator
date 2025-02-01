package main

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

func main() {
	port, targetUrl := parseArgs()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { requestHandler(targetUrl, w, r) })

	log.Printf("Starting server at port %d and forwarding to %s", port, targetUrl)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		log.Println("Error creating new request:", err)
	}
}

func parseArgs() (int, string) {
	if len(os.Args) < 3 {
		log.Fatalf("Usage: tls-terminator <port> <targetHost>")
	}

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	targetUrl := "https://" + os.Args[2]
	_, err = url.ParseRequestURI(targetUrl)
	if err != nil {
		log.Fatalf("Invalid target URL (%s): %v", targetUrl, err)
	}

	return port, targetUrl
}
