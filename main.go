package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
)

const usageMessage = "Usage: tls-exterminator <port>:<host> (or set the CONFIG environment variable)"

func main() {
	config, ok := os.LookupEnv("CONFIG")
	if !ok && len(os.Args) == 2 {
		config = os.Args[1]
	}

	if config == "" {
		fmt.Println(usageMessage)
		log.Fatal("No config provided")
	}

	log.Fatal(run(config))
}

func run(config string) error {
	parts := strings.Split(config, ":")
	if len(parts) != 2 {
		return fmt.Errorf("invalid config: %s", config)
	}

	port, err := strconv.ParseUint(parts[0], 10, 32)
	if err != nil {
		return fmt.Errorf("invalid port %s in config %s", parts[0], config)
	}

	host := parts[1]
	targetUrl := fmt.Sprintf("https://%s", host)

	log.Printf("Starting server at port %d and forwarding to %s", port, targetUrl)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { requestHandler(targetUrl, w, r) })
	return http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}

func requestHandler(targetUrl string, w http.ResponseWriter, r *http.Request) {

	httpsRequest, err := http.NewRequest(r.Method, targetUrl, nil)

	if err != nil {
		log.Println("Error creating new request:", err)
		return
	}

	httpsRequest.URL.Path = r.URL.Path
	httpsRequest.URL.RawQuery = r.URL.RawQuery
	httpsRequest.Header = r.Header
	httpsRequest.Body = r.Body

	// Create a new client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Perform the request
	resp, err := client.Do(httpsRequest)
	if err != nil {
		log.Println("Error creating new request:", err)
		return
	}

	// Set the headers
	for k, v := range resp.Header {
		w.Header()[k] = v
	}
	w.WriteHeader(resp.StatusCode)

	// Copy the body
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println("Error copying response body:", err)
	}
}
