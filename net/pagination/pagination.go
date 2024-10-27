package pagination

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"

	errPkg "github.com/tmeisel/glib/error"
)

var (
	ErrInvalidRequest = errors.New("invalid request")
)

type Request struct {
	Next string `json:"next"`
}

func FromRequest(r *http.Request) (*Request, error) {
	params := r.URL.Query()

	if next, ok := params["pagination.next"]; ok {
		return &Request{Next: next[0]}, nil
	}

	return nil, nil
}

type Response struct {
	More bool    `json:"more"`
	Next *string `json:"next,omitempty"`
}

// Request returns a *Request corresponding
// to the given Response. It returns nil,
// if Response.More is false
func (r Response) Request() *Request {
	if r.More == false {
		return nil
	}

	return &Request{Next: *r.Next}
}

type Pagination interface {
	isPagination() bool
}

type Container struct {
	Type  PaginationType  `json:"type"`
	Inner json.RawMessage `json:"inner"`
}

type Params map[string]interface{}

type LimitAndOffset struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

func (lo LimitAndOffset) isPagination() bool {
	return true
}

// ToResponse encodes the LimitAndOffset as a Response. It will
// always set Response.More to true. If there are no more results,
// it is sufficient to just return &Response{More: false}
func (lo LimitAndOffset) ToResponse() (*Response, error) {
	innerJs, err := json.Marshal(lo)
	if err != nil {
		return nil, errPkg.NewInternalMsg(err, "failed to marshal inner")
	}

	container := Container{
		Type:  TypeLO,
		Inner: innerJs,
	}

	js, err := json.Marshal(container)
	if err != nil {
		return nil, errPkg.NewInternalMsg(err, "failed to marshal pagination container")
	}

	next := base64.URLEncoding.EncodeToString(js)

	return &Response{
		More: true,
		Next: &next,
	}, nil
}

type Token struct {
	Token string `json:"token"`
}

func (t Token) isPagination() bool {
	return true
}

type PaginationType string

const (
	TypeToken = "token"
	TypeLO    = "lo"
)

func (r Request) Decode() (Pagination, error) {
	content := base64.URLEncoding.EncodeToString([]byte(r.Next))

	var c Container
	if err := json.Unmarshal([]byte(content), &c); err != nil {
		return nil, ErrInvalidRequest
	}

	switch c.Type {
	case TypeToken:
		var token Token
		if err := json.Unmarshal(c.Inner, &token); err != nil {
			return nil, ErrInvalidRequest
		}

		return token, nil
	case TypeLO:
		var lo LimitAndOffset
		if err := json.Unmarshal(c.Inner, &lo); err != nil {
			return nil, ErrInvalidRequest
		}

		return lo, nil
	}

	return nil, ErrInvalidRequest
}

func (r Request) AsTokenType() (*Token, error) {
	decoded, err := r.Decode()
	if err != nil {
		return nil, err
	}

	return decoded.(*Token), nil
}

func (r Request) AsLOType() (*LimitAndOffset, error) {
	decoded, err := r.Decode()
	if err != nil {
		return nil, err
	}

	return decoded.(*LimitAndOffset), nil
}
