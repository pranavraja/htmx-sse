package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
)

const CookieName = "JSESSIONID" // a trap of epic proportions

type authenticator struct {
	Secret []byte
}

func (a authenticator) key(token *jwt.Token) (interface{}, error) {
	if alg := token.Method.Alg(); alg != "HS256" {
		return nil, fmt.Errorf("unsupported algorithm %s", alg)
	}
	return a.Secret, nil
}

func (a authenticator) Username(r *http.Request) (string, error) {
	if jwToken, err := r.Cookie(CookieName); err == nil {
		var claims struct {
			jwt.StandardClaims
			Name string
		}
		if _, err := jwt.ParseWithClaims(jwToken.Value, &claims, a.key); err != nil {
			return "", err
		}
		return claims.Name, nil
	}
	return "", nil
}

func (a authenticator) RestrictedHandler(h http.HandlerFunc, username string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		name, err := a.Username(r)
		if err != nil {
			http.Error(w, "sorry, bad token "+err.Error(), http.StatusForbidden)
			return
		}
		if name != "admin" {
			http.Error(w, "sorry, you are not admin", http.StatusForbidden)
			return
		}
		h.ServeHTTP(w, r)
	}
}
