package controllers

import (
	"context"
	"net/http"

	"github.com/AkifhanIlgaz/wit-api/firebase"
)

type UidMiddleware struct {
	AuthService *firebase.AuthService
}

func (umw UidMiddleware) SetUid(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		idToken := r.Header.Get("idToken")

		uid, err := umw.AuthService.GetUidByIdToken(idToken)
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}
		ctx := context.WithValue(r.Context(), "uid", uid)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}
