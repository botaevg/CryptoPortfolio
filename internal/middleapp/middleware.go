package middleapp

import (
	"CryptoPortfolio/internal/servise"
	"context"
	"errors"
	"github.com/dgrijalva/jwt-go/v4"
	"log"
	"net/http"
	"strings"
)

type Auth struct {
	secretkey string
}

func NewAuth(key string) *Auth {
	return &Auth{
		secretkey: key,
	}
}

type Username string

func (a Auth) AuthHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		openURL := []string{"/api/register", "/api/login"}
		path := r.URL.Path
		for _, v := range openURL {
			if v == path {
				next.ServeHTTP(w, r)
				return
			}
		}
		headerAuth := r.Header.Get("Authorization")
		log.Print(headerAuth)
		headerParts := strings.Split(headerAuth, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			http.Error(w, errors.New("unauthorized").Error(), http.StatusUnauthorized)
			return
		}
		tokenClaims := &servise.Claims{}
		token, err := jwt.ParseWithClaims(headerParts[1], tokenClaims, func(token *jwt.Token) (interface{}, error) {
			return []byte(a.secretkey), nil
		})
		if err != nil {
			http.Error(w, errors.New("token error").Error(), http.StatusUnauthorized)
			return
		}
		if !token.Valid {
			http.Error(w, errors.New("token disabled").Error(), http.StatusUnauthorized)
			return
		}
		log.Print(tokenClaims.Username)
		ctx := context.WithValue(r.Context(), Username("username"), tokenClaims.Username)
		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
