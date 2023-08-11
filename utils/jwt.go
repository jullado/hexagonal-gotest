package utils

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var secret = "hexagonal-gotest"

type TokenDataModel struct {
	UserId   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func SignJWT(payload TokenDataModel) (tokenString string, err error) {
	claims := TokenDataModel{
		UserId:   payload.UserId,
		Username: payload.Username,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(1 * time.Hour)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, &claims)

	tokenString, err = token.SignedString([]byte(secret))
	if err != nil {
		return tokenString, errors.New("error sign token")
	}

	return
}

func ValidateJWT(tokenString string) (tokenData TokenDataModel, err error) {

	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		// secret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(secret), nil
	})
	if err != nil {
		return tokenData, err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !(ok && token.Valid) {
		return tokenData, errors.New("error validate token")
	}
	if _, ok := claims["user_id"].(string); !ok {
		return tokenData, errors.New("error validate token")
	}
	if _, ok := claims["username"].(string); !ok {
		return tokenData, errors.New("error validate token")
	}

	dataBytes, err := json.Marshal(claims)
	if err != nil {
		return tokenData, err
	}

	err = json.Unmarshal(dataBytes, &tokenData)
	if err != nil {
		return tokenData, err
	}

	return
}
