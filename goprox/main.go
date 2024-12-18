package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
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
	})

	fmt.Println("Starting server at port 5000")
	if err := http.ListenAndServe(":5000", nil); err != nil {
		log.Println("Error creating new request:", err)
	}
}

func httpsRequest(r *http.Request) (*http.Response, error) {
	// Create a new request using http
	r2, err := http.NewRequest("GET", "https://nabeel.dev", nil)
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
