# Mock Redis

Provide a mock Redis instance and optionally a mock Redis client for testing purposes. This library is built so you can
mock Redis instances using real Redis containers. You will need to have Docker running on your local machine or within
your CI environment.

This library is built on top of [testcontainers](https://testcontainers.com/).

## Usage

Creating a mock instance for creating a customer connection.

```golang
func TestXXX(t *testing.T) {
	ctx := context.Background()

	mock, err := mockredis.NewInstance(ctx, t)

	if err != nil {
		t.Fatalf("creating the instance: %v", err)
		return
	}

	// close the mock
	defer mock.Close(ctx)

	// ... my test code
}
```

Creating a mock redis client for interacting with Redis.

```golang
func TestXXX(t *testing.T) {
	ctx := context.Background()

	mock, err := mockredis.NewClient(ctx, t)

	if err != nil {
		t.Fatalf("creating the client: %v", err)
		return
	}

	// close the mock
	defer mock.Close(ctx)

    redisClient := mock.Client()

	// ... my test code
}
```
