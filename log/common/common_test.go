package common

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/tmeisel/glib/log/fields"
)

func TestUniqueFields(t *testing.T) {
	type testCase struct {
		format         string
		args           []interface{}
		expectedMsg    string
		expectedFields []fields.Field
	}

	for name, tc := range map[string]testCase{
		"msg only": {
			format:      "Here Comes the Hotstepper",
			expectedMsg: "Here Comes the Hotstepper",
		},
		"single arg": {
			format:      "Here Comes the %s",
			args:        []interface{}{"Hotstepper"},
			expectedMsg: "Here Comes the Hotstepper",
		},
		"two args": {
			format:      "Here Comes the %s number %d",
			args:        []interface{}{"Hotstepper", 3},
			expectedMsg: "Here Comes the Hotstepper number 3",
		},
		"args and fields": {
			format:         "Here Comes the %s",
			args:           []interface{}{"Hotstepper", fields.Int("number", 3)},
			expectedMsg:    "Here Comes the Hotstepper",
			expectedFields: []fields.Field{fields.Int("number", 3)},
		},
		"only fields": {
			format:         "Here Comes the Hotstepper",
			args:           []interface{}{fields.Int("number", 3), fields.Error(errors.New("out of Hotsteppers"))},
			expectedMsg:    "Here Comes the Hotstepper",
			expectedFields: []fields.Field{fields.Int("number", 3), fields.Error(errors.New("out of Hotsteppers"))},
		},
	} {
		t.Run(name, func(t *testing.T) {
			output, fieldsOut := ProcessFormatted(tc.format, tc.args...)

			assert.Equal(t, tc.expectedMsg, output)
			assert.Equal(t, tc.expectedFields, fieldsOut)
		})
	}
}
