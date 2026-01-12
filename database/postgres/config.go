package postgres

import (
	"encoding/json"
	"fmt"

	"github.com/tmeisel/glib/exec/backoff"
)

type Config struct {
	Host     string            `envconfig:"HOST" default:"localhost" desc:"Database host"`
	Port     uint              `envconfig:"PORT" default:"5432" desc:"Database port"`
	Database string            `envconfig:"DATABASE" default:"postgres" desc:"Database name"`
	User     string            `envconfig:"USER" default:"postgres" desc:"Database username"`
	Password Password          `envconfig:"PASSWORD" default:"postgres" desc:"Database password"`
	Params   map[string]string `envconfig:"PARAMS" default:"sslmode:require" desc:"Database parameters (as key:value secondKey:secondValue)"`
}

func (c Config) DSN() string {
	var dsn string
	if c.User != "" {
		dsn += c.User
	}

	if c.Password != "" {
		dsn += fmt.Sprintf(":%s", c.Password)
	}

	if dsn != "" {
		dsn += "@"
	}

	dsn += fmt.Sprintf("%s:%d", c.Host, c.Port)

	if c.Database != "" {
		dsn += fmt.Sprintf("/%s", c.Database)
	}

	var params string
	for key, value := range c.Params {
		params += fmt.Sprintf("%s=%s", key, value)
	}

	if params != "" {
		dsn += fmt.Sprintf("?%s", params)
	}

	return fmt.Sprintf("postgres://%s", dsn)
}

type Password string

func (p Password) MarshalJSON() ([]byte, error) {
	replace := make([]rune, len(p))
	for i := range p {
		replace[i] = '*'
	}

	return json.Marshal(replace)
}

type RetryConfig struct {
	Backoff *backoff.Backoff
	Options []backoff.OptionFn
}
