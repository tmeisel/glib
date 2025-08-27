package strutils

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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
		"Negative length": {
			Pos:    0,
			Len:    -5,
			Substr: "World",
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.Substr, SubString(input, tc.Pos, tc.Len))
		})
	}
}

func TestRandom(t *testing.T) {
	type testCase struct {
		expectedLength   int
		expectedAlphabet []string
	}

	for name, tc := range map[string]testCase{
		"lower case": {
			expectedLength: 10,
			expectedAlphabet: []string{
				AlphabetLCChars,
			},
		},
		"combined": {
			expectedLength: 30,
			expectedAlphabet: []string{
				AlphabetLCChars,
				AlphabetNumbers,
			},
		},
	} {
		t.Run(name, func(t *testing.T) {
			output, err := Random(tc.expectedLength, tc.expectedAlphabet...)
			require.NoError(t, err)

			assert.Equal(t, tc.expectedLength, len(output))

		charLoop:
			for _, char := range output {
				for _, alphabet := range tc.expectedAlphabet {
					for _, alpha := range alphabet {
						if char == alpha {
							continue charLoop
						}
					}
				}

				t.Errorf("unexpected character '%c'", char)
			}
		})
	}
}

func TestMustRandom(t *testing.T) {
	assert.NotPanics(t, func() {
		out := MustRandom(10, AlphabetNumbers)
		assert.Equal(t, 10, len(out))
	})
}

func TestInSlice(t *testing.T) {
	const input = "needle"

	type testCase struct {
		Haystack       []string
		ExpectedResult bool
	}

	for name, tc := range map[string]testCase{
		"empty haystack": {
			Haystack:       []string{},
			ExpectedResult: false,
		},
		"single haystack": {
			Haystack:       []string{input},
			ExpectedResult: true,
		},
		"default": {
			Haystack:       []string{"something else", input},
			ExpectedResult: true,
		},
	} {
		t.Run(name, func(t *testing.T) {
			assert.Equal(t, tc.ExpectedResult, InSlice(input, tc.Haystack))
		})
	}
}

func TestInSliceIgnoreCase(t *testing.T) {
	const defaultInput = "Needle"

	type testCase struct {
		Input          *string
		Haystack       []string
		ExpectedResult bool
	}

	for name, tc := range map[string]testCase{
		"empty haystack": {
			Haystack:       []string{},
			ExpectedResult: false,
		},
		"same case": {
			Haystack:       []string{defaultInput},
			ExpectedResult: true,
		},
		"different case": {
			Haystack:       []string{"something else", strings.ToUpper(defaultInput)},
			ExpectedResult: true,
		},
		"empty string": {
			Input:          Ptr(""),
			Haystack:       []string{defaultInput},
			ExpectedResult: false,
		},
	} {
		t.Run(name, func(t *testing.T) {
			input := defaultInput
			if tc.Input != nil {
				input = *tc.Input
			}

			assert.Equal(t, tc.ExpectedResult, InSliceIgnoreCase(input, tc.Haystack))
		})
	}
}
