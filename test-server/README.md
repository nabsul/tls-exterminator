# Test Server

Since tls-exterminator is hard-coded to forward to an https endpoint, we can't test against plain http.
This server uses self-signed certs and serves a response that includes useful debugging data.

## Creating Self-Signed Certs

```sh
openssl req -x509 -nodes -days 3650 -newkey rsa:2048 -keyout server.key -out server.crt -config san.cnf -extensions v3_req
```