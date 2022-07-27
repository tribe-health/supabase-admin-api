package api

import (
	"fmt"
	"net/http"

	"github.com/dgrijalva/jwt-go"
)

type ErrorHandlingWrapper func(w http.ResponseWriter, r *http.Request) error

type Role = string

const (
	SupabaseAdmin Role = "supabase_admin"
	Service       Role = "service_role"
)

func (h ErrorHandlingWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		handleError(err, w, r)
	}
}

func (a *API) RoleValidatingAuthHandler(roleName string) func(next http.Handler) http.Handler {
	wrapping_fn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			tokenString := r.Header.Get("apikey")
			tokenVerifier(w, r, tokenString, a, roleName, next)
		}
		return http.HandlerFunc(fn)
	}
	return wrapping_fn
}

func (a *API) BasicAuthValidatingHandler(roleName string) func(next http.Handler) http.Handler {
	wrapping_fn := func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			_, tokenString, ok := r.BasicAuth()
			if !ok {
				if err := sendJSON(w, http.StatusUnauthorized, "basic auth not provided"); err != nil {
					handleError(err, w, r)
					return
				}
			}
			tokenVerifier(w, r, tokenString, a, roleName, next)
		}
		return http.HandlerFunc(fn)
	}
	return wrapping_fn
}

func tokenVerifier(w http.ResponseWriter, r *http.Request, tokenString string, a *API, roleName string, next http.Handler) {
	if tokenString == "" {
		if err := sendJSON(w, http.StatusUnauthorized, ""); err != nil {
			handleError(err, w, r)
			return
		}
		return
	}

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(a.config.JwtSecret), nil
	})
	if err != nil {
		if err := sendJSON(w, http.StatusUnauthorized, err); err != nil {
			handleError(err, w, r)
			return
		}
		return
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		role, ok := claims["role"]
		if ok && role == roleName {
			// successful authentication
			next.ServeHTTP(w, r)
			return
		} else {
			if err := sendJSON(
				w,
				http.StatusForbidden,
				"this token does not have a valid claim over the correct role",
			); err != nil {
				handleError(err, w, r)
				return
			}
			return
		}
	} else {
		if err := sendJSON(w, http.StatusForbidden, err); err != nil {
			handleError(err, w, r)
			return
		}
		return
	}
}
