package api

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/go-chi/chi"
)

func newRouter() *router {
	return &router{chi.NewRouter()}
}

type router struct {
	chi chi.Router
}

func (r *router) Route(pattern string, fn func(*router)) {
	r.chi.Route(pattern, func(c chi.Router) {
		fn(&router{c})
	})
}

func (r *router) Get(pattern string, fn apiHandler) {
	r.chi.Get(pattern, handler(fn))
}
func (r *router) Post(pattern string, fn apiHandler) {
	r.chi.Post(pattern, handler(fn))
}
func (r *router) Put(pattern string, fn apiHandler) {
	r.chi.Put(pattern, handler(fn))
}
func (r *router) Delete(pattern string, fn apiHandler) {
	r.chi.Delete(pattern, handler(fn))
}

func (r *router) With(fn middlewareHandler) *router {
	c := r.chi.With(middleware(fn))
	return &router{c}
}

func (r *router) WithBypass(fn func(next http.Handler) http.Handler) *router {
	c := r.chi.With(fn)
	return &router{c}
}

func (r *router) Use(fn middlewareHandler) {
	r.chi.Use(middleware(fn))
}
func (r *router) UseBypass(fn func(next http.Handler) http.Handler) {
	r.chi.Use(fn)
}

func (r *router) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	r.chi.ServeHTTP(w, req)
}

type apiHandler func(w http.ResponseWriter, r *http.Request) error

func handler(fn apiHandler) http.HandlerFunc {
	return fn.serve
}

func (h apiHandler) serve(w http.ResponseWriter, r *http.Request) {
	if err := h(w, r); err != nil {
		handleError(err, w, r)
	}
}

type middlewareHandler func(w http.ResponseWriter, r *http.Request) (context.Context, error)

func (m middlewareHandler) handler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Printf("[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI)

		// get jwt token from api header
		tokenString := r.Header.Get("apikey")

		jwtSecret := os.Getenv("JWT_SECRET")

		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Don't forget to validate the alg is what you expect:
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
			return []byte(jwtSecret), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			if claims["role"] == "supabase_admin" {
				// successful authentication
				m.serve(next, w, r)
			} else {
				fmt.Printf("[%s] %s %s %d %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI, http.StatusForbidden, "this token does not have a valid claim over the correct role")
				sendJSON(w, http.StatusForbidden, "this token does not have a valid claim over the correct role")
			}
		} else {
			fmt.Printf("[%s] %s %s %d %s\n", time.Now().Format("2006-01-02 15:04:05"), r.Method, r.RequestURI, http.StatusForbidden, err)
			sendJSON(w, http.StatusForbidden, err)
		}
	})
}

func (m middlewareHandler) serve(next http.Handler, w http.ResponseWriter, r *http.Request) {
	ctx, err := m(w, r)
	if err != nil {
		handleError(err, w, r)
		return
	}
	if ctx != nil {
		r = r.WithContext(ctx)
	}
	next.ServeHTTP(w, r)
}

func middleware(fn middlewareHandler) func(http.Handler) http.Handler {
	return fn.handler
}
