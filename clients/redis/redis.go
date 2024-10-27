package redis

import (
	"github.com/go-redis/redis"

	errPkg "github.com/tmeisel/glib/error"
)

type Config struct {
	Addresses []string `envconfig:"ADDRESS" default:"localhost:6379"`
	Database  int      `envconfig:"DATABASE" default:"0"`
}

type Redis struct {
	client redis.Cmdable
}

// New returns a new Redis client, but does not open a connection.
// To test if the client is able to connect to redis, call
// Redis.Ping
func New(config Config) Redis {
	var client redis.Cmdable

	client = redis.NewUniversalClient(&redis.UniversalOptions{
		Addrs: config.Addresses,
		DB:    config.Database,
	})

	return Redis{client: client}
}

func (r Redis) Ping() error {
	if err := r.client.Ping().Err(); err != nil {
		return errPkg.NewInternal(err)
	}

	return nil
}
