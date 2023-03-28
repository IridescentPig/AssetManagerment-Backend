package utils

import (
	"asset-management/app/define"
	"errors"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const secretTokenSalt = "AssetManagement-BinaryAbstract"

var (
	ErrTokenInvalid = errors.New("token invalid")
	ErrTokenExpire  = errors.New("token expire")
)

const (
	ATokenExpiredDuration = time.Hour
)

func CreateToken(userInfo define.UserBasicInfo) (token string, err error) {
	nowTime := time.Now()
	expiredTime := nowTime.Add(ATokenExpiredDuration)
	stdClaims := jwt.StandardClaims{
		IssuedAt:  nowTime.Unix(),
		NotBefore: nowTime.Unix(),
		ExpiresAt: expiredTime.Unix(),
	}
	claims := define.UserClaims{
		UserBasicInfo:  userInfo,
		StandardClaims: stdClaims,
	}
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

func IsTokenInvalidError(err error) bool {
	return err == ErrTokenInvalid
}

func IsTokenExpiredError(err error) bool {
	return err == ErrTokenExpire
}
