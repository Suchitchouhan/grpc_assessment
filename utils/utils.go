package utils

import (
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const secretkey string = "jn)*hur7x$59tg!lrzosa_c#em)u2yelv%8%*v_j^36ymw"

func GenerateJWT(username string, Password string) string {
	var mySigningKey = []byte(secretkey)
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["username"] = username
	claims["password"] = Password
	tokenString, err := token.SignedString(mySigningKey)
	if err != nil {
		fmt.Printf("Something Went Wrong: %s", err.Error())
		return ""
	}
	return tokenString
}

func DecodeJwt(tokenStr string) (jwt.MapClaims, bool) {
	hmacSecret := []byte(secretkey)
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		return hmacSecret, nil
	})
	if err != nil {
		return nil, false
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims, true
	} else {
		fmt.Printf("Invalid JWT Token")
		return nil, false
	}
}
