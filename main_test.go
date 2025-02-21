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

func TestInvalidConfig(t *testing.T) {
	err := run("")
	require.Error(t, err)
	require.Contains(t, err.Error(), "invalid config")
}

type testRequestData struct {
	name    string
	method  string
	path    string
	query   string
	headers map[string][]string
	body    string
}

var testRequests = []testRequestData{
	{"Basic GET", "GET", "/", "", nil, ""},
	{"GET with path", "GET", "/abc", "", nil, ""},
	{"GET with query params", "GET", "/abc", "q=1,s=2", nil, ""},
	{"POST", "POST", "/", "", nil, "121212"},
}

func TestAll(test *testing.T) {
	for t, h := range portToHost {
		for _, r := range testRequests {
			test.Run(fmt.Sprintf("PORT %s %s", t, r.name), func(test *testing.T) {
				testRequest(test, t, h, r)
			})
		}
	}
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

var portToHost = map[string]string{
	"tls1:5000": "host1",
	"tls2:5001": "host2",
}

func testRequest(t *testing.T, target string, host string, r testRequestData) {
	url := fmt.Sprintf("http://%s", target)
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

	require.Equal(t, host, result.Server)
	require.Equal(t, host, result.Host)
	require.Equal(t, r.method, result.Method)
	require.Equal(t, strings.Trim(r.path, "/"), strings.Trim(result.Path, "/"))
	require.Equal(t, headers, result.Headers)
	require.Equal(t, r.body, result.Body)
	require.Equal(t, r.query, result.Query)
}
