package pagination

import (
	"net/http"
	"net/url"
	"strconv"
	"testing"

	"github.com/tmeisel/glib/utils/strutils"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFromRequest(t *testing.T) {
	type testCase struct {
		Next   *string
		Limit  *uint64
		Offset *uint64

		ExpectedType   PaginationType
		ExpectedLimit  uint64
		ExpectedOffset uint64
		ExpectedError  error
	}

	for name, tc := range map[string]testCase{
		"limit only": {
			Limit:          uintPtr(20),
			ExpectedType:   TypeLimitOffset,
			ExpectedLimit:  20,
			ExpectedOffset: 0,
		},

		"offset only": {
			Offset:         uintPtr(20),
			ExpectedType:   TypeLimitOffset,
			ExpectedLimit:  0,
			ExpectedOffset: 20,
		},

		"limit and offset": {
			Limit:          uintPtr(20),
			Offset:         uintPtr(40),
			ExpectedType:   TypeLimitOffset,
			ExpectedLimit:  20,
			ExpectedOffset: 40,
		},

		"continuation token": {
			Next:         strutils.Ptr("gimme more"),
			ExpectedType: TypeToken,
		},

		"err invalid offset": {
			Limit:         uintPtr(20),
			Offset:        uintPtr(50),
			ExpectedError: ErrInvalidOffset,
		},
		"err invalid mix": {
			Next:          strutils.Ptr("abc"),
			Limit:         uintPtr(20),
			ExpectedError: ErrInvalidMix,
		},
	} {
		t.Run(name, func(t *testing.T) {
			vars := make(url.Values)

			if tc.Limit != nil {
				vars.Set("pagination.limit", strconv.FormatUint(*tc.Limit, 10))
			}

			if tc.Offset != nil {
				vars.Set("pagination.offset", strconv.FormatUint(*tc.Offset, 10))
			}

			if tc.Next != nil {
				next, err := Token{Token: *tc.Next}.ToResponse()
				if err != nil {
					t.Fatal(err)
				}

				vars.Set("pagination.next", *next.Next)
			}

			r := http.Request{URL: &url.URL{RawQuery: vars.Encode()}}

			req, err := FromRequest(&r)
			if tc.ExpectedError != nil {
				require.Error(t, err)
				assert.Equal(t, tc.ExpectedError, err)
				return
			}

			require.NoError(t, err)
			require.NotNil(t, req)

			switch tc.ExpectedType {
			case TypeToken:
				tokenPagination, err := req.AsTokenPagination()
				require.NoError(t, err)
				assert.Equal(t, *tc.Next, tokenPagination.Token)
				return
			case TypeLimitOffset:
				lo, err := req.AsLimitOffsetPagination()
				require.NoError(t, err)
				assert.Equal(t, tc.ExpectedLimit, lo.Limit)
				assert.Equal(t, tc.ExpectedOffset, lo.Offset)
			}

		})
	}
}

func uintPtr(i uint64) *uint64 {
	return &i
}
