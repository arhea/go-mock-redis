package mockredis

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/redis/go-redis/v9"
)

// Client wraps the Redis client and the underlying instance for testing. It provides helper methods for
// accessing the client, closing the client, and closing the instance.
type Client struct {
	t        *testing.T
	instance *Instance
	client   *redis.Client
}

// Client returns the underlying Redis client.
func (r *Client) Client() *redis.Client {
	r.t.Helper()

	return r.client
}

// Close closes the underlying Redis client and the instance.
func (r *Client) Close(ctx context.Context) {
	r.t.Helper()

	if err := r.client.Close(); err != nil {
		r.t.Logf("closing redis client: %v", err)
	}

	r.instance.Close(ctx)
}

// NewClient creates a new Redis client that connects to an underlying instance for testing.
func NewClient(ctx context.Context, t *testing.T) (*Client, error) {
	return NewClientWithConfig(ctx, t, &redis.Options{})
}

// NewClientWithConfig creates a new Redis client that connects to an underlying instance for testing. This method
// will override the `Addr` property with the address of the underlying container.
func NewClientWithConfig(ctx context.Context, t *testing.T, config *redis.Options) (*Client, error) {
	t.Helper()

	instance, err := NewInstance(ctx, t)

	if err != nil {
		return nil, fmt.Errorf("creating the instance: %v", err)
	}

	redisPort, err := instance.Port(ctx)

	if err != nil {
		return nil, fmt.Errorf("getting the mapped port of the instance: %v", err)
	}

	// nolint:errcheck,gosec
	if err := os.Setenv("REDIS_HOST", "localhost"); err != nil {
		return nil, fmt.Errorf("setting the REDIS_HOST environment variable: %v", err)
	}

	if err := os.Setenv("REDIS_PORT", redisPort.Port()); err != nil {
		return nil, fmt.Errorf("setting the REDIS_PORT environment variable: %v", err)
	}

	if err := os.Setenv("REDIS_ADDRESS", fmt.Sprintf("localhost:%s", redisPort.Port())); err != nil {
		return nil, fmt.Errorf("setting the REDIS_ADDRESS environment variable: %v", err)
	}

	// set the address on the config
	config.Addr = fmt.Sprintf("localhost:%s", redisPort.Port())

	redisClient := redis.NewClient(config)

	client := &Client{
		t:        t,
		instance: instance,
		client:   redisClient,
	}

	return client, nil
}
