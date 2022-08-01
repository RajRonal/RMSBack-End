package middleware

import (
	"Rms/database/helper"
	"net/http"
)

func SubAdminMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx := helper.GetContextData(r)
		if ctx == nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		signedUserRole := ctx.Role
		if signedUserRole != "sub-admin" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
