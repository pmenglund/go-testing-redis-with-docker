package main

import (
	"log"
	"net"
	"net/url"
	"os"
	"testing"

	"github.com/go-redis/redis"
	dockertest "gopkg.in/ory-am/dockertest.v3"
)

var client *redis.Client

func TestMain(m *testing.M) {
	// use a test wrapper, as os.Exit ignores defer, so we can't automatically
	// call `pool.Purge(resource)`
	os.Exit(testWrapper(m))
}

func testWrapper(m *testing.M) int {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	resource, err := pool.Run("redis", "latest", nil)
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}
	defer pool.Purge(resource)

	// if run with docker-machine the hostname needs to be set
	u, err := url.Parse(pool.Client.Endpoint())
	if err != nil {
		log.Fatalf("Could not parse endpoint: %s", pool.Client.Endpoint())
	}

	if err := pool.Retry(func() error {
		client = redis.NewClient(&redis.Options{
			Addr:     net.JoinHostPort(u.Hostname(), resource.GetPort("6379/tcp")),
			Password: "", // no password set
			DB:       0,  // use default DB
		})

		ping := client.Ping()
		return ping.Err()
	}); err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	return m.Run()
}

func TestGetNext(t *testing.T) {
	expected := getNext(client)
	if expected != 1 {
		t.Errorf("got %d but expected 1", expected)
	}
}
