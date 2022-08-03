package middleware

import "C"
import (
	"Rms/claims"
	"Rms/database/helper"
	"Rms/handlers"
	"Rms/models"
	"context"
	"fmt"
	"github.com/sirupsen/logrus"
	"net/http"

	"github.com/golang-jwt/jwt"
)

func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tkn := r.Header.Get("Authorization")
		claim := &claims.MapClaims{}
		token, err := jwt.ParseWithClaims(tkn, claim, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return handlers.JwtKey, nil
		})
		if err != nil {
			logrus.Error(err)
		}
		if token.Valid {
			ctx := context.WithValue(r.Context(), models.ClaimKey, claim)
			Context, _ := ctx.Value(models.ClaimKey).(*claims.MapClaims)
			isSession, errs := helper.SessionExist(Context.SessionID)
			if errs != nil {
				logrus.Error("Session Exist : Session does not exist")
				return
			}
			if !isSession {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			logrus.Error(err)
			w.WriteHeader(http.StatusUnauthorized)
			_, err := w.Write([]byte("Unauthorized"))
			if err != nil {
				return
			}
		}
	})
}
