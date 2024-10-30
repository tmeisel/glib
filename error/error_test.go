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
