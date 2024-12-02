package ctx

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIdentity(t *testing.T) {
	myIdentity := "me"

	ctx := WithIdentity(context.Background(), myIdentity)

	assert.Equal(t, myIdentity, GetIdentity(ctx).(string))
}
