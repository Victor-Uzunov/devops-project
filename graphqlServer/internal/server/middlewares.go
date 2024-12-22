package server

import (
	"context"
	"encoding/json"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/constants"
	jwts "github.com/Victor-Uzunov/devops-project/todoservice/pkg/jwt"
	"github.com/Victor-Uzunov/devops-project/todoservice/pkg/log"
	"github.com/golang-jwt/jwt/v4"
	"net/http"
)

func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.C(r.Context()).Info("JWTMiddleware")
		cookie, err := r.Cookie("access_token")
		if err != nil || cookie.Value == "" {
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusForbidden)
			_ = json.NewEncoder(w).Encode(map[string]string{
				"error": "Forbidden",
			})
			return
		}
		claims := &jwts.Claims{}
		token, _ := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return nil, nil
		})
		log.C(r.Context()).Debugf("claims: %v", claims)
		log.C(r.Context()).Debugf("token: %v", token)

		ctx := context.WithValue(r.Context(), constants.TokenCtxKey, cookie.Value)
		ctx = context.WithValue(ctx, "user", claims)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
