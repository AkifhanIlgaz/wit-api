package controllers

import (
	"net/http"

	"github.com/AkifhanIlgaz/wit-api/ctx"
	"github.com/AkifhanIlgaz/wit-api/firebase"
)

type UidMiddleware struct {
	AuthService *firebase.AuthService
}

func (umw UidMiddleware) SetUid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("Authorization")

		uid, err := umw.AuthService.GetUidByIdToken(idToken)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		r = r.WithContext(ctx.WithUid(r.Context(), uid))

		next.ServeHTTP(w, r)
	})
}
