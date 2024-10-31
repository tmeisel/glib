package response

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/tmeisel/glib/utils/strutils"

	"github.com/stretchr/testify/assert"

	dbPkg "github.com/tmeisel/glib/database"
	errPkg "github.com/tmeisel/glib/error"
)

var (
	defaultErr = errors.New("default error")
)

func TestMain(m *testing.M) {
	os.Exit(m.Run())
}

func TestWriteJSON(t *testing.T) {
	type testCase struct {
		Status int
		Body   map[string]interface{}
	}

	for name, tc := range map[string]testCase{
		http.StatusText(http.StatusOK): {
			Status: http.StatusOK,
			Body: map[string]interface{}{
				"key": "value",
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()
			WriteJson(rec, tc.Status, tc.Body)

			assert.Equal(t, tc.Status, rec.Code)

			var r response
			if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
				t.Fatal(err)
			}

			assert.True(t, r.Success)
			assert.Nil(t, r.Error)
			assert.Nil(t, r.Pagination)
			assert.Equal(t, tc.Body, r.Content)
		})
	}
}

func TestWriteError(t *testing.T) {
	type testCase struct {
		Error          error
		ExpectedStatus int
		ExpectedCode   int
	}

	for name, tc := range map[string]testCase{
		"unknown error": {
			Error:          defaultErr,
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedCode:   int(errPkg.CodeInternal),
		},
		"internal pkg err": {
			Error:          errPkg.NewInternal(defaultErr),
			ExpectedStatus: http.StatusInternalServerError,
			ExpectedCode:   int(errPkg.CodeInternal),
		},
		"user pkg err": {
			Error:          errPkg.NewUser(defaultErr),
			ExpectedStatus: http.StatusBadRequest,
			ExpectedCode:   int(errPkg.CodeUser),
		},
		"db err": {
			Error:          dbPkg.NewDuplicateKeyError(nil, strutils.Ptr("id")),
			ExpectedStatus: http.StatusConflict,
			ExpectedCode:   int(errPkg.CodeDuplicateKey),
		},
	} {
		t.Run(name, func(t *testing.T) {
			rec := httptest.NewRecorder()

			WriteError(rec, tc.Error)

			if pkgErr, ok := tc.Error.(*errPkg.Error); ok {
				assert.Equal(t, tc.ExpectedStatus, pkgErr.GetStatus())
			}

			assert.Equal(t, tc.ExpectedStatus, rec.Code)

			var r response
			if err := json.Unmarshal(rec.Body.Bytes(), &r); err != nil {
				t.Fatal(err)
			}

			assert.False(t, r.Success)
			require.NotNil(t, r.Error)
			assert.Equal(t, tc.ExpectedCode, r.Error.Code, rec.Body.String())
		})
	}
}
