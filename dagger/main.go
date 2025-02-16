// A generated module for TlsExterminator functions
//
// This module has been generated via dagger init and serves as a reference to
// basic module structure as you get started with Dagger.
//
// Two functions have been pre-created. You can modify, delete, or add to them,
// as needed. They demonstrate usage of arguments and return types using simple
// echo and grep commands. The functions can be called from the dagger CLI or
// from one of the SDKs.
//
// The first line in this comment block is a short description line and the
// rest is a long description with more detail on the module's purpose or usage,
// if appropriate. All modules should have a short description.

package main

import (
	"context"
	"dagger/tls-exterminator/internal/dagger"
)

type TlsExterminator struct {
	source *dagger.Directory
}

func New(
	ctx context.Context,
	// +defaultPath="./"
	src *dagger.Directory,
) *TlsExterminator {
	return &TlsExterminator{
		source: src,
	}
}

// Builds TLS Exterminator
func (m *TlsExterminator) BuildDockerProd(ctx context.Context) *dagger.Container {
	return m.source.DockerBuild()
}

// Builds TLS Exterminator with cert for testing
func (m *TlsExterminator) BuildDockerTest(ctx context.Context) *dagger.Container {
	return m.BuildDockerProd(ctx).
		WithFile("/etc/ssl/certs/localhost.crt", m.source.File("test-server/localhost.crt"))
}
