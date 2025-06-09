package main

import (
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
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
	targetUrl, err := url.Parse(fmt.Sprintf("https://%s", host))
	if err != nil {
		return fmt.Errorf("invalid host %s in config %s", host, config)
	}

	proxy := httputil.ReverseProxy{
		Rewrite: func(req *httputil.ProxyRequest) {
			req.SetURL(targetUrl)
		},
	}
	http.HandleFunc("/", proxy.ServeHTTP)

	log.Printf("Starting server at port %d and forwarding to %s", port, targetUrl)
	return http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", port), nil)
}
