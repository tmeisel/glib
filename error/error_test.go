package error

import (
	"errors"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
)

var myErrorTest = New(http.StatusBadRequest, "hello", nil)

func TestErrorIs(t *testing.T) {
	assert.True(t, errors.Is(myErrorTest, returnMyError()))
	assert.True(t, returnMyError() == myErrorTest)
}

func returnMyError() error {
	return myErrorTest
}

func TestIsDuplicateKeyErr(t *testing.T) {
	type testCase struct {
		Input    error
		Expected bool
	}

	for name, tc := range map[string]testCase{
		"other error": {
			Input:    NewUser(nil),
			Expected: false,
		},
		"duplicate key error": {
			Input:    New(CodeDuplicateKey, "conflict", nil),
			Expected: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, IsDuplicateKeyErr(tc.Input))
		})
	}

}
