package api

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/middleware"
	"asset-management/utils"
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Init(r *gin.Engine) {
	group := r.Group("/user")

	group.POST("/register", utils.Handler(UserApi.UserRegister))
	group.POST("/login", utils.Handler(UserApi.UserLogin))
	group.GET("/logout", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.UserLogout))
	group.POST("", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.UserCreate))
	group.PATCH("/:username", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.ResetContent))
	group.GET("/:username/lock", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.LockUser))
	group.GET("/:username/unlock", utils.Handler(middleware.JWTMiddleware()), utils.Handler(UserApi.UnlockUser))
	dao.InitForTest()
}

func TestUser(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	Init(r)

	user_register := define.UserRegisterReq{
		UserName: "test",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		body := url.Values{}
		body.Set("userName", user_register.UserName)
		body.Set("password", user_register.Password)
		req, err := http.NewRequest(http.MethodPost, "/user/register", strings.NewReader(body.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "x-www-form-urlencoded")

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		res = httptest.NewRecorder()
		bodyData, err := json.Marshal(user_register)
		if err != nil {
			log.Fatal(err)
		}
		body := bytes.NewReader(bodyData)

		req, err := http.NewRequest(http.MethodPost, "/user/register", body)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	{
		res = httptest.NewRecorder()
		bodyData, err := json.Marshal(user_register)
		if err != nil {
			log.Fatal(err)
		}
		body := bytes.NewReader(bodyData)

		req, err := http.NewRequest(http.MethodPost, "/user/register", body)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	user_login := define.UserLoginReq{
		UserName: "test",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		res = httptest.NewRecorder()
		body := url.Values{}
		body.Set("userName", user_login.UserName)
		body.Set("password", user_login.Password)
		req, err := http.NewRequest(http.MethodPost, "/user/login", strings.NewReader(body.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "x-www-form-urlencoded")

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		bodyData, err := json.Marshal(user_login)
		if err != nil {
			log.Fatal(err)
		}
		body := bytes.NewReader(bodyData)

		req, err := http.NewRequest(http.MethodPost, "/user/login", body)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	b, err := io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}

	{
		data := map[string]interface{}{}
		json.Unmarshal(b, &data)
		// fmt.Println(data["data"])
		user := data["data"].(map[string]interface{})
		token := user["token"].(string)
		req, err := http.NewRequest(http.MethodGet, "/user/logout", nil)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", token)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}
}

func TestAdmin(t *testing.T) {
	res := httptest.NewRecorder()
	_, r := gin.CreateTestContext(res)
	Init(r)

	admin := model.User{
		UserName:     "admin",
		Password:     utils.CreateMD5("21232f297a57a5a743894a0e4a801fc3"),
		SystemSuper:  true,
		EntityID:     nil,
		DepartmentID: nil,
		Ban:          false,
	}
	dao.UserDao.Create(admin)

	userLogin := define.UserLoginReq{
		UserName: "admin",
		Password: "21232f297a57a5a743894a0e4a801fc3",
	}

	{
		bodyData, err := json.Marshal(userLogin)
		if err != nil {
			log.Fatal(err)
		}
		body := bytes.NewReader(bodyData)

		req, err := http.NewRequest(http.MethodPost, "/user/login", body)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

	b, err := io.ReadAll(res.Result().Body)
	if err != nil {
		log.Fatal(err)
	}

	data := map[string]interface{}{}
	json.Unmarshal(b, &data)
	// fmt.Println(data["data"])
	user := data["data"].(map[string]interface{})
	token := user["token"].(string)

	userCreate := define.UserRegisterReq{
		UserName: "test2",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	{
		res = httptest.NewRecorder()
		body := url.Values{}
		body.Set("userName", userCreate.UserName)
		body.Set("password", userCreate.Password)
		req, err := http.NewRequest(http.MethodPost, "/user", strings.NewReader(body.Encode()))
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Authorization", token)
		req.Header.Set("Content-Type", "x-www-form-urlencoded")

		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusBadRequest, res.Result().StatusCode, "response failed")
	}

	{
		bodyData, err := json.Marshal(userCreate)
		if err != nil {
			log.Fatal(err)
		}
		body := bytes.NewReader(bodyData)

		req, err := http.NewRequest(http.MethodPost, "/user", body)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")

		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusUnauthorized, res.Result().StatusCode, "response failed")
	}

	{
		bodyData, err := json.Marshal(userCreate)
		if err != nil {
			log.Fatal(err)
		}
		body := bytes.NewReader(bodyData)

		req, err := http.NewRequest(http.MethodPost, "/user", body)
		if err != nil {
			log.Fatal(err)
		}
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", token)
		res = httptest.NewRecorder()
		r.ServeHTTP(res, req)

		assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	}

}
