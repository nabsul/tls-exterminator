package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAll(t *testing.T) {
	for p := range portToHost {
		testReq(t, p, "GET", "/", "", nil, "")
		testReq(t, p, "GET", "/abc", "", nil, "")
		testReq(t, p, "GET", "/abc", "q=1,s=2", nil, "")
		testReq(t, p, "POST", "/", "", nil, "121212")
	}
}

func testReq(t *testing.T, port int, method, path, query string, headers map[string][]string, body string) {
	r := requestData{
		port:    port,
		method:  method,
		path:    path,
		query:   query,
		headers: headers,
		body:    body,
	}

	testRequest(t, r)
}

type requestData struct {
	port    int
	method  string
	path    string
	query   string
	headers map[string][]string
	body    string
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

var portToHost = map[int]string{
	5000: "host1",
	5001: "host2",
}

func testRequest(t *testing.T, r requestData) {
	url := fmt.Sprintf("http://localhost:%d", r.port)
	req, err := http.NewRequest(r.method, url, io.NopCloser(strings.NewReader(r.body)))
	if err != nil {
		t.Errorf("Error creating request: %v", err)
	}

	req.Header = r.headers
	req.Method = r.method
	req.URL.Path = r.path
	req.URL.RawQuery = r.query

	client := &http.Client{}
	resp, err := client.Do(req)
	require.NoError(t, err)

	data, err := io.ReadAll(resp.Body)
	require.NoError(t, err)

	result := Response{}
	err = json.Unmarshal(data, &result)
	require.NoError(t, err)

	headers := make(map[string][]string)
	headers["Accept-Encoding"] = []string{"gzip"}
	headers["User-Agent"] = []string{"Go-http-client/1.1"}
	if r.headers != nil {
		for k, v := range r.headers {
			headers[k] = v
		}
	}

	require.Equal(t, portToHost[r.port], result.Server)
	require.Equal(t, portToHost[r.port], result.Host)
	require.Equal(t, r.method, result.Method)
	require.Equal(t, strings.Trim(r.path, "/"), strings.Trim(result.Path, "/"))
	require.Equal(t, headers, result.Headers)
	require.Equal(t, r.body, result.Body)
	require.Equal(t, r.query, result.Query)
}
