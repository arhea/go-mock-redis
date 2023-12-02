package mockredis_test

import (
	"context"
	"testing"

	mockredis "github.com/arhea/go-mock-redis"
)

func TestInstance(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock, err := mockredis.NewInstance(ctx, t)

	if err != nil {
		t.Fatalf("creating the instance: %v", err)
		return
	}

	// close the mock
	defer mock.Close(ctx)

	port, err := mock.Port(ctx)

	if err != nil {
		t.Fatalf("getting the mapped port of the instance: %v", err)
		return
	}

	if port.Port() == "" {
		t.Fatalf("port should not be empty")
		return
	}
}
