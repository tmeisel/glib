package response

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/tmeisel/glib/net/pagination"

	errPkg "github.com/tmeisel/glib/error"
)

type response struct {
	Success    bool                 `json:"success"`
	Error      *Error               `json:"error,omitempty"`
	Pagination *pagination.Response `json:"pagination,omitempty"`
	Content    any                  `json:"content,omitempty"`
}

type Error struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func WriteError(w http.ResponseWriter, err error) {
	status := http.StatusInternalServerError
	if pkgErr, ok := err.(*errPkg.Error); ok {
		status = pkgErr.GetStatus()
	}

	WriteErrorStatus(w, status, err)
}

func WriteErrorStatus(w http.ResponseWriter, status int, err error) {
	code := status
	if pkgErr, ok := err.(*errPkg.Error); ok {
		code = int(pkgErr.GetCode())
	}

	// from e.g. http status 500 => code 50000
	if code < 1000 {
		code = code * 100
	}

	writeJson(w, status, response{Success: false, Error: &Error{
		Code:    code,
		Message: err.Error(),
	}})
}

func WriteJson(w http.ResponseWriter, status int, v any) {
	writeJson(w, status, response{Success: true, Content: v})
}

func WritePaginated(w http.ResponseWriter, v any, pagination *pagination.Response) {
	writeJson(w, http.StatusOK, response{Success: true, Pagination: pagination, Content: v})
}

// ETag sets an ETag header with the given value val
func ETag(w http.ResponseWriter, weak bool, val string) {
	if weak {
		val = fmt.Sprintf(`W/"%s"`, val)
	}

	w.Header().Set("ETag", val)
}

// LastModified sets a Last-Modified Header using the given timestamp t
func LastModified(w http.ResponseWriter, t time.Time) {
	w.Header().Set("Last-Modified", t.Format(http.TimeFormat))
}

func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to write json to ResponseWriter: %+e", err)
	}
}
