package middleware

import (
	"context"
	"errors"
	"github.com/thorgnir-go-study/yp-diploma/internal/pkg/auth"
	"net/http"
)

var (
	ClaimsCtxKey = &contextKey{"JWTClaims"}
)

var (
	ErrClaimsNotPresent = errors.New("there are no claims in provided context")
	ErrClaimsInvalid    = errors.New("invalid cast")
)

func AuthMiddleware(j *auth.JwtWrapper) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		hfn := func(w http.ResponseWriter, r *http.Request) {
			cookie, err := r.Cookie("jwt")
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			if cookie.Value == "" {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}

			claims, err := j.ValidateToken(cookie.Value)
			if err != nil {
				http.Error(w, http.StatusText(http.StatusUnauthorized), http.StatusUnauthorized)
				return
			}
			ctx := newContext(r.Context(), *claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		}
		return http.HandlerFunc(hfn)
	}
}

func GetClaimsFromContext(ctx context.Context) (*auth.CustomClaims, error) {
	rawClaims := ctx.Value(ClaimsCtxKey)
	if rawClaims == nil {
		return nil, ErrClaimsNotPresent
	}
	claims, ok := rawClaims.(auth.CustomClaims)
	if !ok {
		return nil, ErrClaimsInvalid
	}
	return &claims, nil
}

func newContext(ctx context.Context, claims auth.CustomClaims) context.Context {
	ctx = context.WithValue(ctx, ClaimsCtxKey, claims)
	return ctx
}

// Утащено из go-chi/jwtauth
// contextKey is a value for use with context.WithValue. It's used as
// a pointer so it fits in an interface{} without allocation. This technique
// for defining context keys was copied from Go 1.7's new use of context in net/http.
type contextKey struct {
	name string
}
