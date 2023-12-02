package mockredis

import (
	"context"
	"testing"
	"time"

	"github.com/cenkalti/backoff/v4"
	"github.com/docker/go-connections/nat"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

// Instance represents the underlying container that is running the mock redis instance.
type Instance struct {
	t         *testing.T
	container testcontainers.Container
}

// Port returns the mapped port of the underlying container.
func (r *Instance) Port(ctx context.Context) (nat.Port, error) {
	r.t.Helper()

	return r.container.MappedPort(ctx, "6379")
}

// Close terminates the underlying container.
func (r *Instance) Close(ctx context.Context) {
	r.t.Helper()

	if err := r.container.Terminate(ctx); err != nil {
		r.t.Logf("error terminating redis emulator: %v", err)
	}
}

// NewInstance creates a new Redis container. It will exponentially backoff until the container is ready to accept
// connections so that you can handle throttling within CI environments
func NewInstance(ctx context.Context, t *testing.T) (*Instance, error) {
	// mark this a test helper function
	t.Helper()

	// configure the backoff
	cfg := backoff.NewExponentialBackOff()
	cfg.InitialInterval = time.Second * 1
	cfg.MaxElapsedTime = time.Minute * 5
	policy := backoff.WithContext(cfg, ctx)

	// configure our retry logic when connecting to docker, this is helpful when running tests in parallel and on ci
	operation := backoff.OperationWithData[testcontainers.Container](func() (testcontainers.Container, error) {
		return testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
			ContainerRequest: testcontainers.ContainerRequest{
				Image:        "docker.io/redis:6",
				ExposedPorts: []string{"6379/tcp"},
				WaitingFor:   wait.ForLog("* Ready to accept connections"),
			},
			Started: true,
			Reuse:   false,
			Logger:  testcontainers.TestLogger(t),
		})
	})

	// create the redis emulator container
	redisEmulator, err := backoff.RetryWithData(operation, policy)

	if err != nil {
		return nil, err
	}

	// create the mock instance type
	cntr := &Instance{
		t:         t,
		container: redisEmulator,
	}

	return cntr, nil
}
