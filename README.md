# TLS Exterminator

Credits: Special thanks to @gerhard for working through the design and testing of all of this!

A small service that exposes an HTTP endpoint and reverse proxies all requests to a domain over HTTPS.
The goal is to use this program as a sidecar to work around a limitation where the client only does HTTP.

## Background

In [this episode of The Changelog](https://changelog.com/friends/73#transcript-519)
there is a discussion of a limitation (in Varnish I think)
where TLS termination is not supported.
This got me thinking: I can solve that!
This repo contains a few experiments and the solution I've settled on to work around that limitation.

Dotnet/C# is my favorite programming language, so I started there.
But then I realized that Rust or Go have two distinct advantages:

- Smaller binaries and memory footprints
- More widely accepted by the open source community

All things equal, I probably would have gone with Rust, but I decided to use Go because:

- I have much more experience with Go than Rust
- Everything needed is built into Go's standard library, while a Rust implementation would require some third-party libraries

## Build

If you just want to build the binary on your local machine

```sh
go build cmd
```

## Usage

```sh
./tls-exterminator [listen-port]:[target-domain]
```

## Test Locally

