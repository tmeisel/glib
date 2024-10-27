package response

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/tmeisel/glib/pagination"

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
	if pkgErr, ok := err.(errPkg.Error); ok {
		status = pkgErr.GetStatus()
	}

	WriteErrorStatus(w, status, err)
}

func WriteErrorStatus(w http.ResponseWriter, status int, err error) {
	code := status
	if pkgErr, ok := err.(errPkg.Error); ok {
		code = pkgErr.GetCode()
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

func writeJson(w http.ResponseWriter, status int, v any) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)

	if err := json.NewEncoder(w).Encode(v); err != nil {
		log.Printf("failed to write json to ResponseWriter: %+e", err)
	}
}
