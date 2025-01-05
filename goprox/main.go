package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
)

var targetUrl string

func main() {

	port, err := strconv.Atoi(os.Args[1])
	if err != nil {
		log.Fatalf("Invalid port number: %v", err)
	}

	targetUrl = "https://" + os.Args[2]

	http.HandleFunc("/", requestHandler)

	portStr := strconv.Itoa(port)
	fmt.Println("Starting server at port " + portStr + " and forwarding to " + targetUrl)
	if err := http.ListenAndServe(":"+portStr, nil); err != nil {
		log.Println("Error creating new request:", err)
	}
}

func requestHandler(w http.ResponseWriter, r *http.Request) {
	resp, err := httpsRequest(r)
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

func httpsRequest(r *http.Request) (*http.Response, error) {
	// Create a new request using http
	r2, err := http.NewRequest("GET", targetUrl, nil)
	if err != nil {
		return nil, err
	}

	r2.URL.Path = r.URL.Path
	copyHeaders(r.Header, r2.Header)
	r2.Body = r.Body

	// Create a new client
	client := &http.Client{}

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
