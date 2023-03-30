package api

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/utils"
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func Init(r *gin.Engine) {
	group := r.Group("/user")

	group.POST("/register", utils.Handler(UserApi.UserRegister))
	group.POST("/login", utils.Handler(UserApi.UserLogin))

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

	// body := url.Values{}
	// body.Set("userName", user_register.UserName)
	// body.Set("password", user_register.Password)
	// req, err := http.NewRequest(http.MethodPost, "/user/register", strings.NewReader(body.Encode()))
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
	// b, err := io.ReadAll(res.Result().Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// data := map[string]interface{}{}
	// json.Unmarshal(b, &data)
	// fmt.Println(data["error"])

	user_login := define.UserLoginReq{
		UserName: "test",
		Password: "e10adc3949ba59abbe56e057f20f883e",
	}

	// body = url.Values{}
	// body.Set("userName", user_login.UserName)
	// body.Set("password", user_login.Password)
	// req, err = http.NewRequest(http.MethodPost, "/user/login", strings.NewReader(body.Encode()))
	bodyData, err = json.Marshal(user_login)
	if err != nil {
		log.Fatal(err)
	}
	body = bytes.NewReader(bodyData)

	req, err = http.NewRequest(http.MethodPost, "/user/login", body)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")

	res = httptest.NewRecorder()
	r.ServeHTTP(res, req)

	assert.Equal(t, http.StatusOK, res.Result().StatusCode, "response failed")
	// b, err := io.ReadAll(res.Result().Body)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// data := map[string]interface{}{}
	// json.Unmarshal(b, &data)
	// fmt.Println(data["data"])
}
