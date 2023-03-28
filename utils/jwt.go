package utils

import (
	"errors"
	"fmt"

	"github.com/dgrijalva/jwt-go"
)

const secretTokenSalt = "AssetManagement-BinaryAbstract"

var (
	ErrTokenInvalid = errors.New("token invalid")
	ErrTokenExpire  = errors.New("token expire")
)

func CreateToken(claims jwt.MapClaims) (token string, err error) {
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err = tokenObj.SignedString([]byte(secretTokenSalt))
	return
}

func ParseToken(token string) (claims jwt.MapClaims, err error) {
	tokenObj, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secretTokenSalt), nil
	})

	if err != nil {
		validationErr, ok := err.(jwt.ValidationError)
		if ok && (validationErr.Errors&(jwt.ValidationErrorExpired) != 0) {
			err = ErrTokenExpire
			return
		}
	}

	var isOK bool
	claims, isOK = tokenObj.Claims.(jwt.MapClaims)

	if !isOK {
		err = ErrTokenInvalid
		return
	}

	return
}
