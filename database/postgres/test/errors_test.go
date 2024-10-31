package test

import (
	"context"
	"embed"
	"log"
	"os"
	"testing"
	"time"

	"github.com/jackc/pgx/v5/pgconn"

	"github.com/tmeisel/glib/database"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	postgresPkg "github.com/tmeisel/glib/database/postgres"
	"github.com/tmeisel/glib/database/postgres/docker"
	errPkg "github.com/tmeisel/glib/error"
)

var (
	//go:embed test_migrations
	migrations embed.FS
	container  *docker.Postgres
)

func TestMain(m *testing.M) {
	ctx := context.Background()

	var err error
	container, err = docker.NewPostgres(ctx, docker.V16, time.Minute, time.Second*10)
	if err != nil {
		log.Printf("failed to start postgres container: %v", err)
		os.Exit(1)
	}

	if err := container.RunMigrations(ctx, migrations, "test_migrations"); err != nil {
		log.Printf("failed to run migrations: %v", err)
		os.Exit(1)
	}

	os.Exit(m.Run())
}

func TestProcessError(t *testing.T) {
	ctx := context.Background()
	conn := container.GetPool()

	const id = "abc"

	_, err := conn.Exec(ctx, "INSERT INTO users (id) VALUES ($1);", id)
	require.NoError(t, err)

	_, err = conn.Exec(ctx, "INSERT INTO users (id) VALUES ($1);", id)
	require.Error(t, err)

	err = postgresPkg.ProcessError(err)
	require.Error(t, err)

	pkgErr, ok := err.(*errPkg.Error)
	require.True(t, ok)
	assert.Equal(t, errPkg.CodeDuplicateKey, pkgErr.GetCode())

}

func TestIsDuplicateKeyError(t *testing.T) {
	type testCase struct {
		Input    error
		Expected bool
	}

	for name, tc := range map[string]testCase{
		"nil": {
			Input:    nil,
			Expected: false,
		},
		"errPkg": {
			Input:    database.NewDuplicateKeyError(nil, nil),
			Expected: true,
		},
		"pgErr": {
			Input:    &pgconn.PgError{Code: postgresPkg.CodeDuplicateKey},
			Expected: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, postgresPkg.IsDuplicateKeyError(tc.Input))
		})
	}
}
