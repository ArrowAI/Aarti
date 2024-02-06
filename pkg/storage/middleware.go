package storage

import (
	"net/http"

	"aarti/pkg/utils/logger"

	"github.com/gorilla/mux"

	"aarti/pkg/auth"
)

type MiddlewareFunc = func(repoVar string) mux.MiddlewareFunc

func Middleware(ar Repository) MiddlewareFunc {
	return func(repoVar string) mux.MiddlewareFunc {
		return func(next http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				ctx := auth.Context(r.Context(), r)
				name := mux.Vars(r)[repoVar]
				if name == "" {
					n := Options(ctx).repo
					if n == "" {
						http.Error(w, "missing repository name", http.StatusBadRequest)
						return
					}
					name = n
				}
				ctx = logger.Set(ctx, logger.C(ctx).WithField("repo", name).WithField("type", ar.Name()))
				s, err := NewStorage(ctx, name, ar)
				if err != nil {
					Error(w, err)
					return
				}
				defer s.Close()
				next.ServeHTTP(w, r.WithContext(Context(ctx, s)))
			})
		}
	}
}
