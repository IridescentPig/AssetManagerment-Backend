package middleware

import (
	"asset-management/app/define"
	"asset-management/utils"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func HelloFunc(ctx *utils.Context) {
	ctx.JSON(200, gin.H{
		"message": "hello",
	})
}

func TestCors(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(Cors())
	r.GET("/hello", utils.Handler(HelloFunc))

	req, err := http.NewRequest(http.MethodGet, "/hello", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("origin", "test")
	res := httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode, "cors error")

	req, err = http.NewRequest(http.MethodOptions, "/hello", nil)
	if err != nil {
		log.Fatal(err)
	}
	res = httptest.NewRecorder()
	r.ServeHTTP(res, req)
	assert.Equal(t, http.StatusOK, res.Result().StatusCode, "cors error")
}

func TestJwt(t *testing.T) {
	w := httptest.NewRecorder()
	_, r := gin.CreateTestContext(w)

	r.Use(utils.Handler(JWTMiddleware()))
	r.GET("/hello", utils.Handler(HelloFunc))

	{
		req, err := http.NewRequest(http.MethodGet, "/hello", nil)
		if err != nil {
			log.Fatal(err)
		}
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "jwt middleware error")
	}
	{
		req, err := http.NewRequest(http.MethodGet, "/hello", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", "123456")
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "jwt middleware error")
	}
	{
		nowTime := time.Now()
		m, _ := time.ParseDuration("-2m")
		issuedTime := nowTime.Add(m)
		m, _ = time.ParseDuration("-1m")
		expiredTime := nowTime.Add(m)
		stdClaims := jwt.StandardClaims{
			IssuedAt:  issuedTime.Unix(),
			NotBefore: issuedTime.Unix(),
			ExpiresAt: expiredTime.Unix(),
		}
		tokenObj := jwt.NewWithClaims(jwt.SigningMethodHS256, stdClaims)
		token, err := tokenObj.SignedString([]byte("AssetManagement-BinaryAbstract"))
		assert.Equal(t, nil, err, "jwt error")

		req, err := http.NewRequest(http.MethodGet, "/hello", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", token)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "jwt middleware error")
	}
	{
		token, err := utils.CreateToken(define.UserBasicInfo{
			UserID:          1,
			UserName:        "admin",
			EntitySuper:     true,
			DepartmentSuper: true,
			SystemSuper:     true,
		})
		assert.Equal(t, nil, err, "jwt create error")

		req, err := http.NewRequest(http.MethodGet, "/hello", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", token)
		res := httptest.NewRecorder()
		r.ServeHTTP(res, req)
		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "jwt middleware error")
	}
}
