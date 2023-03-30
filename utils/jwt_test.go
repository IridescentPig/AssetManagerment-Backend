package utils

import (
	"asset-management/app/define"
	"testing"

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

	t.Logf("Token is %s.", token)
}

func TestParseToken(t *testing.T) {
	var err error

	t.Logf("Token is %s.", token)
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
}
