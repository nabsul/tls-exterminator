FROM golang:1.23 AS build

WORKDIR /app
COPY . .
RUN go build main.go

FROM alpine AS certs
RUN apk update && apk add ca-certificates

FROM busybox
COPY --from=certs /etc/ssl/certs /etc/ssl/certs
WORKDIR /app
COPY --from=build /app/main /app/main

CMD ["/app/main"]
