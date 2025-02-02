package main

import (
	"fmt"
	"io"
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

func requestHandler(targetUrl string, w http.ResponseWriter, r *http.Request) {
	resp, err := httpsRequest(targetUrl, r)
	if err != nil {
		log.Println("Error creating new request:", err)
		return
	}

	copyHeaders(resp.Header, w.Header())
	w.WriteHeader(resp.StatusCode)

	// Copy the body
	defer resp.Body.Close()
	_, err = io.Copy(w, resp.Body)
	if err != nil {
		log.Println("Error copying response body:", err)
	}
}

func httpsRequest(targetUrl string, r *http.Request) (*http.Response, error) {
	// Create a new request using http
	r2, err := http.NewRequest("GET", targetUrl, nil)

	if err != nil {
		return nil, err
	}

	r2.URL.Path = r.URL.Path
	copyHeaders(r.Header, r2.Header)
	r2.Body = r.Body

	// Create a new client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Perform the request
	resp, err := client.Do(r2)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

func copyHeaders(from http.Header, to http.Header) {
	for key, values := range from {
		for _, value := range values {
			to.Add(key, value)
		}
	}
}
