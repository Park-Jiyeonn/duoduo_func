package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

type MyClaims struct {
	Username string
	jwt.StandardClaims
}

var key = []byte("Jiyeon_Hyomin_Jiyeon_Hyomin_hhhh")

func GetToken(username string) (string, error) {
	claims := MyClaims{
		username,
		jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour * 24 * 7).Unix(),
			Issuer:    "朴智妍",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(key)
}

func ParseToken(token string) (*MyClaims, error) {
	ret, err := jwt.ParseWithClaims(token, &MyClaims{}, func(token *jwt.Token) (i interface{}, err error) {
		return key, nil
	})
	if err != nil {
		return nil, err
	}
	claims, ok := ret.Claims.(*MyClaims)
	if !ok || !ret.Valid {
		return nil, fmt.Errorf("invalid token")
	}
	return claims, nil
}
