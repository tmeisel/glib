package docker

import (
	"context"
	"embed"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgxpool"

	postgresPkg "github.com/tmeisel/glib/database/postgres"
	dockerPkg "github.com/tmeisel/glib/testing/docker"
)

const (
	repo = "postgres"

	username = "postgres"
	password = "postgres"
	database = "postgres"
)

const (
	V13 = "13-alpine"
	V14 = "14-alpine"
	V15 = "15-alpine"
	V16 = "16-alpine"
)

type Postgres struct {
	*dockerPkg.Container

	dsn string
	pgx *pgxpool.Pool

	exposedPort uint
}

// NewPostgres runs a PostgreSQL server in a docker container. The maxContainerLifetime must
// be set according to the tests it is used for. The waitTime specifies how long you want to
// wait for the client to successfully connect to the postgres server AFTER the container
// started. Hence, downloading the docker image and starting the container may exceed the
// given waitTime.
// If the waitTime is eventually exceeded, an error will be returned.
// version needs to be a valid tag, e.g. (postgres:)13-alpine.
//
// The caller should eventually call Cleanup when tests are finished (usually in TestMain)
func NewPostgres(ctx context.Context, version string, maxContainerLifetime, waitTime time.Duration) (*Postgres, error) {
	t := &Postgres{}

	if err := t.initContainer(version, maxContainerLifetime); err != nil {
		return nil, err
	}

	if err := t.waitForContainer(ctx, waitTime); err != nil {
		return nil, err
	}

	return t, nil
}

func (t *Postgres) GetPool() *pgxpool.Pool {
	return t.pgx
}

func (t *Postgres) RunMigrations(ctx context.Context, migrations embed.FS, path string) error {
	return postgresPkg.Migrate(ctx, migrations, path, t.GetDSN())
}

func (t *Postgres) GetDSN() string {
	if t.dsn == "" {
		t.dsn = fmt.Sprintf(
			"postgres://%s:%s@localhost:%d/%s?sslmode=disable",
			username,
			password,
			t.exposedPort,
			database,
		)
	}

	return t.dsn
}

func (t *Postgres) GetPSQL() squirrel.StatementBuilderType {
	return squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar)
}

// Truncate removes all rows from the given tables in a single statement
func (t *Postgres) Truncate(ctx context.Context, tables ...string) error {
	stmt := fmt.Sprintf("TRUNCATE TABLE %s", strings.Join(tables, ","))

	if _, err := t.pgx.Exec(ctx, stmt); err != nil {
		return err
	}

	return nil
}

func (t *Postgres) Close() error {
	t.pgx.Close()
	return t.Container.Cleanup()
}

func (t *Postgres) initContainer(version string, maxLifetime time.Duration) error {
	var err error

	env := []string{
		fmt.Sprintf("POSTGRES_USERNAME=%s", username),
		fmt.Sprintf("POSTGRES_PASSWORD=%s", password),
		fmt.Sprintf("POSTGRES_DB=%s", database),
	}

	// create the container pool
	t.Container, err = dockerPkg.New(repo, version, maxLifetime, env...)
	if err != nil {
		return err
	}

	portString := t.Container.GetPortString("5432/tcp")

	port, err := strconv.Atoi(portString)
	if err != nil {
		return err
	}

	t.exposedPort = uint(port)

	return err
}

func (t *Postgres) waitForContainer(ctx context.Context, maxWaitTime time.Duration) error {
	timeout := time.NewTimer(maxWaitTime)

	conf := postgresPkg.Config{
		Host:     "localhost",
		Port:     t.exposedPort,
		User:     username,
		Password: password,
		Database: database,
		Params: map[string]string{
			"SSLMode": "disable",
		},
	}

	var err error
	t.pgx, err = pgxpool.New(ctx, conf.DSN())
	if err != nil {
		return err
	}

	for {
		if err := t.pgx.Ping(ctx); err == nil {
			return nil
		}

		retry := time.NewTicker(time.Second).C

		select {
		case <-timeout.C:
			return err
		case <-retry:
			continue
		}
	}
}
