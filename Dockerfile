FROM alpine AS certs
RUN apk update && apk add ca-certificates

FROM golang:1.23 AS build
WORKDIR /app
COPY . .
RUN go build -o tls-exterminator . 

FROM busybox
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
WORKDIR /app
COPY --from=build /app/tls-exterminator /app/tls-exterminator

CMD ["/app/tls-exterminator"]
