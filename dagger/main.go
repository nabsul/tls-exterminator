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

const (
	goBuildTag = "golang:1.23"
	goRunTag   = "busybox:latest"
)

type TlsExterminator struct{}

// Returns a container that echoes whatever string argument is provided
func (m *TlsExterminator) ContainerEcho(stringArg string) *dagger.Container {
	return dag.Container().From("alpine:latest").WithExec([]string{"echo", stringArg})
}

// Returns lines that match a pattern in the files of the provided Directory
func (m *TlsExterminator) GrepDir(ctx context.Context, directoryArg *dagger.Directory, pattern string) (string, error) {
	return dag.Container().
		From("alpine:latest").
		WithMountedDirectory("/mnt", directoryArg).
		WithWorkdir("/mnt").
		WithExec([]string{"grep", "-R", pattern, "."}).
		Stdout(ctx)
}

// Builds the main TLS Exterminator binary
func (m *TlsExterminator) Build(ctx context.Context, src *dagger.Directory) *dagger.Container {
	build := dag.Container().
		From(goBuildTag).
		WithMountedDirectory("/src", src).
		WithExec([]string{"ls"}).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "/tls-extermintator", "."})

	binary := build.File("/tls-extermintator")
	certs := build.Directory("/etc/ssl/certs")

	return dag.Container().
		From(goRunTag).
		WithDirectory("/etc/ssl/certs", certs, dagger.ContainerWithDirectoryOpts{Exclude: []string{"_archive/", ".git/"}}).
		WithFile("/app/tls-extermintator", binary).
		WithWorkdir("/app").
		WithEntrypoint([]string{"/app/tls-extermintator"})
}

// Builds the test server binary
func (m *TlsExterminator) BuildTestServer(ctx context.Context, src *dagger.Directory) *dagger.Container {
	build := dag.Container().
		From(goBuildTag).
		WithMountedDirectory("/src", src).
		WithExec([]string{"ls"}).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "/test-server", "./test-server"})

	binary := build.File("/test-server")
	key := src.File("test-server/server.key")
	cert := src.File("test-server/server.crt")

	return dag.Container().
		From(goRunTag).
		WithFile("/app/test-server", binary).
		WithFile("/app/server.key", key).
		WithFile("/app/server.crt", cert).
		WithWorkdir("/app").
		WithEntrypoint([]string{"/app/test-server"})

}

// Builds the test TLS Exterminator binary
func (m *TlsExterminator) BuildTestTlsExterminator(ctx context.Context, src *dagger.Directory) *dagger.Container {
	cert := src.File("test-server/server.crt")
	return m.Build(ctx, src).
		WithFile("/etc/ssl/certs/server.crt", cert).
		WithEntrypoint([]string{"/app/tls-exterminator"})
}

func (m *TlsExterminator) Test(ctx context.Context, src *dagger.Directory) (string, error) {
	testServer := m.BuildTestServer(ctx, src)
	tls := m.BuildTestTlsExterminator(ctx, src)

	srv1 := testServer.
		WithExposedPort(443).
		WithEnvVariable("HOST", "host1").
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{"/app/test-server"},
		})
	srv2 := testServer.
		WithExposedPort(443).
		WithEnvVariable("HOST", "host2").
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{"/app/test-server"},
		})

	tls1 := tls.
		WithServiceBinding("host1", srv1).
		WithExposedPort(5000).
		WithEnvVariable("CONFIG", "5000:host1").
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{"/app/tls-exterminator"},
		})
	tls2 := tls.
		WithServiceBinding("host2", srv2).
		WithExposedPort(5001).
		WithEnvVariable("CONFIG", "5001:host2").
		AsService(dagger.ContainerAsServiceOpts{
			Args: []string{"/app/tls-exterminator"},
		})

	return m.Build(ctx, src).
		WithServiceBinding("tls1", tls1).
		WithServiceBinding("tls1", tls2).
		WithExec([]string{"go", "test", "./..."}).
		Stdout(ctx)
}
