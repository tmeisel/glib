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

func TestIs(t *testing.T) {
	type testCase struct {
		Error    error
		Code     Code
		Expected bool
	}

	for name, tc := range map[string]testCase{
		"other error": {
			Error:    errors.New("some error"),
			Code:     CodeUser,
			Expected: false,
		},
		"duplicate key error": {
			Error:    New(CodeDuplicateKey, "msg", nil),
			Code:     CodeDuplicateKey,
			Expected: true,
		},
		"not duplicate key error": {
			Error:    New(CodeDuplicateKey, "msg", nil),
			Code:     CodeUser,
			Expected: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Expected, Is(tc.Error, tc.Code))
		})
	}
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
