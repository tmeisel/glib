package pagination

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"

	errPkg "github.com/tmeisel/glib/error"
)

var (
	ErrInvalidRequest = errPkg.NewUserMsg(nil, "invalid request")
	ErrInvalidMix     = errPkg.NewUserMsg(nil, "invalid mix of pagination.next and pagination.limit/pagination.offset")
	ErrInvalidUint    = errPkg.NewUserMsg(nil, "invalid value for either pagination.limit or pagination.offset")
	ErrInvalidOffset  = errPkg.NewUserMsg(nil, "invalid offset. expecting a multiple of limit")
)

type Request struct {
	Next string `json:"next"`
}

func FromRequest(r *http.Request) (*Request, error) {
	params := r.URL.Query()

	if params.Has("pagination.next") && (params.Has("pagination.limit") || params.Has("pagination.offset")) {
		return nil, ErrInvalidMix
	}

	if next, ok := params["pagination.next"]; ok {
		return &Request{Next: next[0]}, nil
	}

	if params.Has("pagination.limit") || params.Has("pagination.offset") {
		var lo LimitAndOffset

		if limit, ok := params["pagination.limit"]; ok {
			lim, err := strconv.ParseUint(limit[0], 10, 64)
			if err != nil {
				return nil, ErrInvalidUint
			}

			lo.Limit = lim
		}

		if offset, ok := params["pagination.offset"]; ok {
			off, err := strconv.ParseUint(offset[0], 10, 64)
			if err != nil {
				return nil, ErrInvalidUint
			}

			lo.Offset = off
		}

		if lo.Offset > 0 && lo.Limit > 0 && lo.Offset%lo.Limit != 0 {
			return nil, ErrInvalidOffset
		}

		return lo.toRequest()
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

type Container struct {
	Type  PaginationType  `json:"type"`
	Inner json.RawMessage `json:"inner"`
}

type Params map[string]interface{}

type LimitAndOffset struct {
	Limit  uint64 `json:"limit"`
	Offset uint64 `json:"offset"`
}

func (lo LimitAndOffset) toRequest() (*Request, error) {
	js, err := json.Marshal(lo)
	if err != nil {
		return nil, err
	}

	container := Container{
		Type:  TypeLimitOffset,
		Inner: js,
	}

	jsContainer, err := json.Marshal(container)
	if err != nil {
		return nil, err
	}

	return &Request{Next: base64.URLEncoding.EncodeToString(jsContainer)}, nil
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
		Type:  TypeLimitOffset,
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

func (t Token) ToResponse() (*Response, error) {
	innerJs, err := json.Marshal(t)
	if err != nil {
		return nil, errPkg.NewInternalMsg(err, "failed to marshal inner")
	}

	container := Container{
		Type:  TypeToken,
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

type PaginationType string

const (
	TypeToken       PaginationType = "token"
	TypeLimitOffset PaginationType = "limitOffset"
)

func (r Request) decodeContainer() (*Container, error) {
	containerJS, err := base64.URLEncoding.DecodeString(r.Next)
	if err != nil {
		return nil, errPkg.NewUserMsg(err, "failed to decode next page")
	}

	var c Container
	if err := json.Unmarshal(containerJS, &c); err != nil {
		return nil, ErrInvalidRequest
	}

	return &c, nil
}

func (r Request) AsTokenPagination() (*Token, error) {
	container, err := r.decodeContainer()
	if err != nil {
		return nil, err
	}

	if container.Type != TypeToken {
		return nil, ErrInvalidRequest
	}

	var t Token
	if err := json.Unmarshal(container.Inner, &t); err != nil {
		return nil, ErrInvalidRequest
	}

	return &t, nil
}

func (r Request) AsLimitOffsetPagination() (*LimitAndOffset, error) {
	container, err := r.decodeContainer()
	if err != nil {
		return nil, err
	}

	if container.Type != TypeLimitOffset {
		return nil, ErrInvalidRequest
	}

	var lo LimitAndOffset
	if err := json.Unmarshal(container.Inner, &lo); err != nil {
		return nil, ErrInvalidRequest
	}

	return &lo, nil
}
