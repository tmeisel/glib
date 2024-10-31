package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5/pgconn"
	"gotest.tools/assert"

	"github.com/tmeisel/glib/database"
)

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
			Input:    &pgconn.PgError{Code: CodeDuplicateKey},
			Expected: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, IsDuplicateKeyError(tc.Input))
		})
	}
}
