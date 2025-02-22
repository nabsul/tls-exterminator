package main

import (
	"context"
	"dagger/tls-exterminator/internal/dagger"
)

const (
	goBuildTag = "golang:1.23"
	goRunTag   = "busybox:latest"
)

type TlsExterminator struct {
	// +private
	Src *dagger.Directory
}

func New(
	// +defaultPath="./"
	// +ignore=["_archive"]
	src *dagger.Directory,
) *TlsExterminator {
	return &TlsExterminator{
		Src: src,
	}
}

// Builds the main TLS Exterminator binary
func (m *TlsExterminator) Build(ctx context.Context) *dagger.Container {
	build := dag.Container().
		From(goBuildTag).
		WithDirectory("/src", m.Src).
		WithExec([]string{"ls"}).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "/tls-exterminator", "."})

	binary := build.File("/tls-exterminator")
	certs := build.Directory("/etc/ssl/certs")

	return dag.Container().
		From(goRunTag).
		WithDirectory("/etc/ssl/certs", certs).
		WithFile("/app/tls-exterminator", binary).
		WithWorkdir("/app").
		WithEntrypoint([]string{"/app/tls-exterminator"})
}

// Builds the test server binary
func (m *TlsExterminator) BuildTestServer(ctx context.Context) *dagger.Container {
	build := dag.Container().
		From(goBuildTag).
		WithDirectory("/src", m.Src).
		WithExec([]string{"ls"}).
		WithWorkdir("/src").
		WithExec([]string{"go", "build", "-o", "/test-server", "./test-server"})

	binary := build.File("/test-server")
	key := m.Src.File("test-server/server.key")
	cert := m.Src.File("test-server/server.crt")

	return dag.Container().
		From(goRunTag).
		WithWorkdir("/app").
		WithFile("/app/test-server", binary).
		WithFile("/app/server.key", key).
		WithFile("/app/server.crt", cert).
		WithEntrypoint([]string{"/app/test-server"})

}

// Builds the test TLS Exterminator binary
func (m *TlsExterminator) BuildTestTlsExterminator(ctx context.Context) *dagger.Container {
	cert := m.Src.File("test-server/server.crt")
	return m.Build(ctx).
		WithFile("/etc/ssl/certs/server.crt", cert).
		WithEntrypoint([]string{"/app/tls-exterminator"})
}

// Runs integration tests
func (m *TlsExterminator) Test(ctx context.Context) (string, error) {
	testServer := m.BuildTestServer(ctx)
	tls := m.BuildTestTlsExterminator(ctx)

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
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true})

	tls2 := tls.
		WithServiceBinding("host2", srv2).
		WithExposedPort(5001).
		WithEnvVariable("CONFIG", "5001:host2").
		AsService(dagger.ContainerAsServiceOpts{UseEntrypoint: true})

	return dag.Container().
		From(goBuildTag).
		WithWorkdir("/src").
		WithDirectory("/src", m.Src).
		WithServiceBinding("tls1", tls1).
		WithServiceBinding("tls2", tls2).
		WithExec([]string{"go", "test", "-v", "."}).
		Stdout(ctx)
}
