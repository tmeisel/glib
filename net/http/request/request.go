package request

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"time"

	errPkg "github.com/tmeisel/glib/error"

	"github.com/gorilla/mux"
)

var (
	ErrNoSuchHeader = errPkg.NewUserMsg(nil, "no such header")
)

func GetRouteParam(r *http.Request, key string) string {
	return mux.Vars(r)[key]
}

// IfNoneMatch returns true, if the given Etag does not equal the requests If-None-Match header
func IfNoneMatch(r *http.Request, ETag string) (bool, error) {
	headerValue := r.Header.Get("If-None-Match")
	if headerValue == "" {
		return false, ErrNoSuchHeader
	}

	return headerValue != ETag, nil
}

// IfMatch returns true, if the given Etag equals the requests If-Match header
func IfMatch(r *http.Request, ETag string) bool {
	return r.Header.Get("If-Match") == ETag
}

// IfModifiedSince returns true, if lastUpdate is greater than
// the time specified in the requests If-Modified-Since header.
func IfModifiedSince(r *http.Request, lastUpdate time.Time) (bool, error) {
	headerValue := r.Header.Get("If-Modified-Since")
	if headerValue == "" {
		return false, ErrNoSuchHeader
	}

	since, err := time.Parse(time.RFC1123, headerValue)
	if err != nil {
		return false, errPkg.NewUserMsg(err, "Could not parse If-Modified-Since header")
	}

	return lastUpdate.After(since), nil
}

// IfUnmodifiedSince returns true, if lastUpdate is not greater than
// the time specified in the requests If-Unmodified-Since header.
func IfUnmodifiedSince(r *http.Request, lastUpdate time.Time) (bool, error) {
	headerValue := r.Header.Get("If-Unmodified-Since")
	if headerValue == "" {
		return false, ErrNoSuchHeader
	}

	since, err := time.Parse(time.RFC1123, headerValue)
	if err != nil {
		return false, errPkg.NewUserMsg(err, "Could not parse If-Unmodified Since header")
	}

	return lastUpdate.After(since) == false, nil
}

func DecodeBody(r *http.Request, dest any) error {
	if err := json.NewDecoder(r.Body).Decode(&dest); err != nil {
		if err == io.EOF {
			return errPkg.New(http.StatusBadRequest, "no request body specified", nil)
		}

		return errPkg.New(http.StatusBadRequest, "invalid request body specified (invalid json)", err)
	}

	return nil
}

const (
	AuthHeader   = "Authorization"
	ApiKeyHeader = "X-API-Key"
	BearerPrefix = "Bearer"
)

// GetAuthToken returns the token specified in the
// Authorization header of the request. That could be
// the actual Authorization header, prefixed with "Bearer "
// or the X-API-Key header
func GetAuthToken(r *http.Request) string {
	header := getAuthHeader(r)
	if header == "" {
		return ""
	}

	return strings.TrimSpace(strings.TrimPrefix(header, BearerPrefix))

}

func getAuthHeader(r *http.Request) string {
	if header := r.Header.Get(AuthHeader); header != "" {
		return header
	}

	return r.Header.Get(ApiKeyHeader)
}
