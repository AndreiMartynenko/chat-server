package root

import (
	"errors"
	"io"
	"os"
	"time"

	jwt "github.com/golang-jwt/jwt/v5"

	"github.com/AndreiMartynenko/common/pkg/closer"
)

type userClaims struct {
	jwt.RegisteredClaims
	Username string `json:"username"`
	Role     string `json:"role"`
}

func readToken() (string, error) {
	file, err := os.Open(filename)
	if err != nil {
		return "", err
	}
	closer.Add(file.Close)

	token, err := io.ReadAll(file)
	if err != nil {
		return "", err
	}

	return string(token), nil
}

func getTokenClaims(token string) (*userClaims, error) {
	t, err := jwt.ParseWithClaims(
		token,
		&userClaims{},
		func(_ *jwt.Token) (interface{}, error) {
			return []byte(""), nil
		})
	if !errors.Is(err, jwt.ErrTokenSignatureInvalid) {
		return nil, err
	}

	claims, ok := t.Claims.(*userClaims)
	if !ok {
		return nil, errors.New("invalid access token")
	}

	return claims, nil
}

func isTokenExpired(claims *userClaims) error {
	expire, err := claims.RegisteredClaims.GetExpirationTime()
	if err != nil {
		return err
	}

	if expire.Before(time.Now()) {
		// Token expired
		return errors.New("access token expired")
	}

	return nil
}
