package common

import (
	"fmt"

	"github.com/tmeisel/glib/log/fields"
)

// ProcessFormatted splits args belonging to the format string from
// additional fields. It then returns the formatted message and
// the fields separately
func ProcessFormatted(format string, args ...interface{}) (string, []fields.Field) {
	var fieldsOut []fields.Field
	for _, arg := range args {
		if field, ok := arg.(fields.Field); ok {
			fieldsOut = append(fieldsOut, field)
		}
	}

	fieldCount := len(fieldsOut)
	if fieldCount == 0 {
		return fmt.Sprintf(format, args...), nil
	}

	firstField := len(args) - fieldCount

	return fmt.Sprintf(format, args[:firstField]...), fieldsOut
}

// JoinUnique returns all fields from input plus all fields from
// more but unique according to their Field.Key
func JoinUnique(input []fields.Field, more ...fields.Field) []fields.Field {
	input = append(input, more...)

	output := make([]fields.Field, 0)
	keys := make(map[string]int)

	for _, field := range input {
		idx, exists := keys[field.Key]
		if exists {
			output[idx] = field
			continue
		}

		keys[field.Key] = len(output)
		output = append(output, field)
	}

	return output
}
