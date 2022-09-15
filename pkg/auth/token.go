package auth

import "github.com/golang-jwt/jwt"

type Token struct {
	UserID int
	Name   string
	Email  string
	Role   string
	Claims *jwt.StandardClaims
}

func (t Token) Valid() error {
	return nil
}

// func CreateJWT() (string, error){
// 	token := jwt.New(jwt.SigningMethodHS256)

// 	claims := token.Claims.(jwt.MapClaims)

// 	claims["exp"] =
// }
