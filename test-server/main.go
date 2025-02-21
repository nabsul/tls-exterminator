package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

const (
	cert = "server.crt"
	key  = "server.key"
)

func main() {
	host, ok := os.LookupEnv("HOST")
	if !ok {
		log.Fatal("HOST environment variable is required")
	}

	log.Fatal(Run(host, cert, key))
}

func Run(host, cert, key string) error {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) { handleRequest(host, w, r) })
	fmt.Printf("Serving host [%s] at port 443\n", host)
	return http.ListenAndServeTLS(":443", cert, key, nil)
}

type Response struct {
	Server  string
	Host    string
	Method  string
	Path    string
	Query   string
	Headers map[string][]string
	Body    string
}

func handleRequest(host string, w http.ResponseWriter, r *http.Request) {

	res, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("Error reading request body:", err)
		res = []byte("Error reading request body")
	}

	resp := Response{
		Server:  host,
		Host:    r.Host,
		Method:  r.Method,
		Path:    r.URL.Path,
		Query:   r.URL.RawQuery,
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
