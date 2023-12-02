package mockredis_test

import (
	"context"
	"testing"

	mockredis "github.com/arhea/go-mock-redis"
	"github.com/stretchr/testify/assert"
)

func TestClient(t *testing.T) {
	t.Parallel()

	ctx := context.Background()

	mock, err := mockredis.NewClient(ctx, t)

	if err != nil {
		t.Fatalf("creating the client: %v", err)
		return
	}

	// close the mock
	defer mock.Close(ctx)

	assert.Nil(t, mock.Client().Ping(ctx).Err(), "ping should work")
	assert.NotNil(t, mock.Client(), "client should not be nil")
}
