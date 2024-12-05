package auth

import (
	"context"
	"net/http"

	ctxPkg "github.com/tmeisel/glib/ctx"

	"github.com/tmeisel/glib/net/http/request"
	"github.com/tmeisel/glib/net/http/response"
)

type AuthMiddleware struct {
	identityFn IdentityFunc
}

type Identity interface {
	Valid() error
}

// IdentityFunc should take a bearerToken and return the
// corresponding identity
type IdentityFunc func(ctx context.Context, bearerToken string) (Identity, error)

func NewAuthMiddleware(identityFunc IdentityFunc) *AuthMiddleware {
	return &AuthMiddleware{
		identityFn: identityFunc,
	}
}

// IdentityMiddleware extracts the bearer token from the request and adds
// the corresponding identity to the request context
func (a *AuthMiddleware) IdentityMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if authToken := request.GetAuthToken(r); authToken == "" {
			identity, err := a.getIdentity(r.Context(), authToken)
			if err != nil {
				response.WriteError(w, err)
				return
			}

			r = r.WithContext(ctxPkg.WithIdentity(r.Context(), identity))
		}

		next.ServeHTTP(w, r)
	})
}

// RequireIdentity calls Identity.Valid and cancels the request,
// if it returns an error
func (a *AuthMiddleware) RequireIdentity(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		identity, err := a.getIdentity(r.Context(), request.GetAuthToken(r))
		if err != nil {
			response.WriteError(w, err)
			return
		}

		if err := identity.Valid(); err != nil {
			response.WriteError(w, err)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (a *AuthMiddleware) getIdentity(ctx context.Context, bearerToken string) (Identity, error) {
	if a.identityFn == nil {
		panic("identityFn is nil")
	}

	return a.identityFn(ctx, bearerToken)
}
