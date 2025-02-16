package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func startServer(port uint64, host string) error {
	targetUrl := fmt.Sprintf("https://%s", host)
	log.Printf("Starting server at port %d and forwarding to %s", port, targetUrl)
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { requestHandler(targetUrl, w, r) })
	return http.ListenAndServe(fmt.Sprintf(":%d", port), mux)
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

func httpsRequest(targetUrl string, httpRequest *http.Request) (*http.Response, error) {
	// Create a new request using http
	httpsRequest, err := http.NewRequest(httpRequest.Method, targetUrl, nil)

	if err != nil {
		return nil, err
	}

	httpsRequest.URL.Path = httpRequest.URL.Path
	copyHeaders(httpRequest.Header, httpsRequest.Header)
	httpsRequest.Body = httpRequest.Body

	// Create a new client
	client := &http.Client{
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		},
	}

	// Perform the request
	return client.Do(httpsRequest)
}

func copyHeaders(from http.Header, to http.Header) {
	for key, values := range from {
		for _, value := range values {
			to.Add(key, value)
		}
	}
}
