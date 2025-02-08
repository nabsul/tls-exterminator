package main

import (
	"fmt"
	"log"
)

const usageMessage = "Usage: tls_exterminator <port>:<host>,<port>:<host>,... (or set the CONFIG environment variable)"

func main() {
	config, err := parseArgs()
	if err != nil {
		fmt.Println(usageMessage)
		log.Fatal(err)
	}

	errorChan := make(chan error)
	for port, targetUrl := range config {
		go func() { errorChan <- startServer(port, targetUrl) }()
	}

	log.Fatal(<-errorChan)
}
