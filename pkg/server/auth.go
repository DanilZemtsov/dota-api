package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"net/http"
	"strings"
	"time"
)

func newJWT(userID string, t time.Duration, key string) (string, error) {

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		ExpiresAt: time.Now().Add(t).UnixNano(),
		Subject:   userID,
	})

	return token.SignedString([]byte(key))
}

func parseJWT(acctoken string, key string) (string, error) {
	token, err := jwt.Parse(acctoken, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("")
		}
		return []byte(key), nil
	})
	if err != nil {
		return "", err
	}
	clams, ok := token.Claims.(jwt.MapClaims)
	if !ok {

	}
	return clams["sub"].(string), nil
}
func authHeader(r *http.Request) (string, error) {
	authHader := r.Header.Get("token")
	token := strings.Replace(authHader, "Bearer ", "", 1)

	userId, err := parseJWT(token, SecretKeyToken)
	if err != nil {
		return "", err
	}
	return userId, nil
}
