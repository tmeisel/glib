package request

import (
	"encoding/json"
	"io"
	"net/http"

	errPkg "github.com/tmeisel/glib/error"

	"github.com/gorilla/mux"
)

func GetRouteParam(r *http.Request, key string) string {
	return mux.Vars(r)[key]
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
