package main

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyCustomClaims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func ParseMyCustomClaims(tokenString string, key []byte) (*MyCustomClaims, error) {
	// 解析JWT令牌
	token, err := jwt.ParseWithClaims(tokenString, &MyCustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}

	// 检查令牌是否有效
	claims, ok := token.Claims.(*MyCustomClaims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func main() {
	// create claims
	claims := MyCustomClaims{
		"my-username",
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24).Unix(),
			Issuer:    "my-issuer",
		},
	}

	// create token with claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	fmt.Println(token)
	// sign the token with a secret key
	signedToken, err := token.SignedString([]byte("my-secret-key"))
	if err != nil {
		panic(err)
	}
	fmt.Println(signedToken)

	fmt.Println(ParseMyCustomClaims(signedToken, []byte("my-secret-key")))
}
