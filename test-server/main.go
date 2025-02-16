package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("HOST environment variable is required")
	}

	cert := "server.crt"
	key := "server.key"

	http.HandleFunc("/", handleRequest)
	fmt.Printf("Serving host [%s] at port 443\n", host)
	log.Fatal(http.ListenAndServeTLS(":443", cert, key, nil))
}

type Response struct {
	Host    string
	Method  string
	Url     string
	Headers map[string][]string
	Body    string
}

func handleRequest(w http.ResponseWriter, r *http.Request) {

	res, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		res = []byte("Error reading request body")
	}

	resp := Response{
		Host:    r.Host,
		Method:  r.Method,
		Url:     r.URL.String(),
		Headers: r.Header,
		Body:    string(res),
	}

	json, err := json.Marshal(resp)
	if err != nil {
		log.Println("Error marshalling response:", err)
		w.WriteHeader(500)
		w.Write([]byte("Error marshalling response"))
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(json)
}
