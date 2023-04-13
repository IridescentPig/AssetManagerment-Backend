package utils

import (
	"asset-management/app/define"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/stretchr/testify/assert"
)

var userInfo define.UserBasicInfo
var claims *define.UserClaims
var token string

func TestCreateToken(t *testing.T) {
	userInfo = define.UserBasicInfo{
		UserID:          10000,
		UserName:        "test",
		EntitySuper:     false,
		DepartmentSuper: false,
		SystemSuper:     false,
	}

	var err error

	token, err = CreateToken(userInfo)
	if err != nil {
		t.Fatalf("Token generation error.")
	}

	// t.Logf("Token is %s.", token)
}

func TestParseToken(t *testing.T) {
	var err error

	// t.Logf("Token is %s.", token)
	claims, err = ParseToken(token)

	if IsTokenExpiredError(err) {
		t.Fatalf("Token has expired.")
	}

	if IsTokenInvalidError(err) {
		t.Fatalf("Token invalid.")
	}

	assert.Equal(t, uint(10000), claims.UserID)
	assert.Equal(t, "test", claims.UserName)
	assert.Equal(t, false, claims.EntitySuper)
	assert.Equal(t, false, claims.DepartmentSuper)
	assert.Equal(t, false, claims.DepartmentSuper)

	nowTime := time.Now()
	expiredTime := nowTime.Add(time.Hour)
	stdClaims := jwt.StandardClaims{
		IssuedAt:  nowTime.Unix(),
		NotBefore: expiredTime.Unix(),
		ExpiresAt: nowTime.Unix(),
	}
	tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
	token, err = tokenObj.SignedString([]byte(secretTokenSalt))
	assert.Equal(t, nil, err, "jwt error")

	claims, err = ParseToken(token)
	assert.Equal(t, ErrTokenInvalid, err, "jwt error")

	nowTime = time.Now()
	m, _ := time.ParseDuration("-2m")
	issuedTime := nowTime.Add(m)
	m, _ = time.ParseDuration("-1m")
	expiredTime = nowTime.Add(m)
	stdClaims = jwt.StandardClaims{
		IssuedAt:  issuedTime.Unix(),
		NotBefore: issuedTime.Unix(),
		ExpiresAt: expiredTime.Unix(),
	}
	tokenObj = jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
	token, err = tokenObj.SignedString([]byte(secretTokenSalt))
	assert.Equal(t, nil, err, "jwt error")

	claims, err = ParseToken(token)
	assert.Equal(t, ErrTokenExpire, err, "jwt error")

	nowTime = time.Now()
	expiredTime = nowTime.Add(time.Hour)
	mapClaims := jwt.MapClaims{
		"id":       1,
		"username": "test",
		"nbf":      nowTime.Unix(),
		"iat":      expiredTime.Unix(),
	}
	tokenObj = jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaims)
	token, err = tokenObj.SignedString([]byte(secretTokenSalt))
	assert.Equal(t, nil, err, "jwt error")

	claims, err = ParseToken(token)
	assert.Equal(t, ErrTokenInvalid, err, "jwt error")
}
