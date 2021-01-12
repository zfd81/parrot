package otto

import (
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

type RockClaims struct {
	Data map[string]interface{}
	jwt.StandardClaims
}

func CreateToken(data map[string]interface{}, secret string, second int) string {
	now := time.Now()
	claims := RockClaims{
		Data: data,
		StandardClaims: jwt.StandardClaims{
			IssuedAt:  now.Unix(),
			ExpiresAt: now.Add(time.Second * time.Duration(second)).Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(secret))
	if err != nil {
		throwException(err.Error())
	}
	return tokenString

}

func ParseToken(tokenString string, secret string) map[string]interface{} {
	token, err := jwt.ParseWithClaims(tokenString, &RockClaims{}, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				throwException("That's not even a token")
			} else if ve.Errors&jwt.ValidationErrorExpired != 0 {
				throwException("Token is expired")
			} else if ve.Errors&jwt.ValidationErrorNotValidYet != 0 {
				throwException("Token not active yet")
			} else {
				throwException("Couldn't handle this token")
			}
		}
	}
	if claims, ok := token.Claims.(*RockClaims); ok && token.Valid {
		return claims.Data
	}
	throwException("Couldn't handle this token")
	return nil
}