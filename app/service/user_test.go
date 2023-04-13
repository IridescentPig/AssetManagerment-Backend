package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/utils"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func InitForTest() {
	dao.InitForTest()
}

func TestUser(t *testing.T) {
	InitForTest()

	err := UserService.CreateUser("admin", "admin")
	if err != nil {
		log.Fatal(err)
	}

	var exist bool
	exist, err = UserService.ExistsUser("admin")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, true, exist, "service error")
	exist, err = UserService.ExistsUser("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, false, exist, "service error")

	var (
		token    string
		userInfo *model.User
	)

	token, userInfo, err = UserService.VerifyPasswordAndGetUser("test", "123456")
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, (*model.User)(nil), userInfo, "service error")
	assert.Equal(t, "", token, "servicve error")
	token, userInfo, err = UserService.VerifyPasswordAndGetUser("admin", "123456")
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, (*model.User)(nil), userInfo, "service error")
	assert.Equal(t, "", token, "servicve error")
	token, userInfo, err = UserService.VerifyPasswordAndGetUser("admin", "admin")
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, uint(1), userInfo.ID, "service error")
	assert.Equal(t, false, token == "", "servicve error")

	err = UserService.ModifyUserIdentity("admin", 1)
	assert.Equal(t, nil, err, "service error")
	_, userInfo, err = UserService.VerifyPasswordAndGetUser("admin", "admin")
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, true, userInfo.DepartmentSuper, "service error")

	err = UserService.ModifyUserPassword("admin", "123456")
	assert.Equal(t, nil, err, "service error")
	_, userInfo, err = UserService.VerifyPasswordAndGetUser("admin", "123456")
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, utils.CreateMD5("123456"), userInfo.Password, "service error")

	err = UserService.ModifyUserBanstate("admin", true)
	assert.Equal(t, nil, err, "service error")
	_, userInfo, err = UserService.VerifyPasswordAndGetUser("admin", "123456")
	assert.Equal(t, nil, err, "service error")
	assert.Equal(t, true, userInfo.Ban, "service error")

	res := httptest.NewRecorder()
	ctx, _ := gin.CreateTestContext(res)
	context := utils.Context{Context: ctx}

	isSuper := UserService.SystemSuper(&context)
	assert.Equal(t, false, isSuper, "service error")
	isSuper = UserService.EntitySuper(&context)
	assert.Equal(t, false, isSuper, "service error")
	isSuper = UserService.DepartmentSuper(&context)
	assert.Equal(t, false, isSuper, "service error")
	username, err := UserService.UserName(&context)
	assert.Equal(t, "", username, "service error")
	assert.Equal(t, "no user vertification info", err.Error(), "service error")

	ctx.Set("user", define.UserBasicInfo{
		UserID:          1,
		UserName:        "admin",
		EntitySuper:     true,
		DepartmentSuper: true,
		SystemSuper:     true,
	})
	isSuper = UserService.SystemSuper(&context)
	assert.Equal(t, true, isSuper, "service error")
	isSuper = UserService.EntitySuper(&context)
	assert.Equal(t, true, isSuper, "service error")
	isSuper = UserService.DepartmentSuper(&context)
	assert.Equal(t, true, isSuper, "service error")
	username, err = UserService.UserName(&context)
	assert.Equal(t, "admin", username, "service error")
	assert.Equal(t, nil, err, "service error")

}
