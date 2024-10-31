package strutils

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSubString(t *testing.T) {
	const input = "Hello World"

	type testCase struct {
		Pos    int
		Len    int
		Substr string
	}

	for name, tc := range map[string]testCase{
		"Hello": {
			Pos:    0,
			Len:    5,
			Substr: "Hello",
		},
		"World": {
			Pos:    6,
			Len:    5,
			Substr: "World",
		},
		"Pos gt Input": {
			Pos:    30,
			Len:    5,
			Substr: "",
		},
		"Pos plus Len gt Input": {
			Pos:    6,
			Len:    15,
			Substr: "World",
		},
		"Len eq 0": {
			Pos:    6,
			Len:    0,
			Substr: "World",
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Substr, SubString(input, tc.Pos, tc.Len))
		})
	}
}
