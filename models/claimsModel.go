package models

import (
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v4"
)

type JwtClaims struct {
	UserId string
	Roles  []string
	jwt.StandardClaims
}

func (claims JwtClaims) Valid() error {
	var now = time.Now().UTC().Unix()
	if claims.VerifyExpiresAt(now, true) {
		return nil
	}
	return fmt.Errorf("Token is invalid")
}
