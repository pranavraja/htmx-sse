package main

import (
	"fmt"
	"os"

	"github.com/golang-jwt/jwt/v4"
)

func main() {
	var claims struct {
		jwt.StandardClaims
		Name string
	}
	claims.Name = os.Args[1]
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, _ := token.SignedString([]byte(os.Getenv("SESSION_SECRET")))
	fmt.Println(signed)
}
