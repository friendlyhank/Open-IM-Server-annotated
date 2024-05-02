package authverify

import "github.com/golang-jwt/jwt/v4"

func Secret(secret string) jwt.Keyfunc {
	return func(token *jwt.Token) (any, error) {
		return []byte(secret), nil
	}
}
