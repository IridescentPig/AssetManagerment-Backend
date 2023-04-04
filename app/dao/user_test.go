package dao

import (
	"asset-management/app/model"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Init() {
	InitForTest()
}

func TestUser(t *testing.T) {
	Init()
	user := model.User{
		UserName: "test",
		Password: "123456",
	}
	err := UserDao.Create(user)
	assert.Equal(t, nil, err, "database error")

	user = model.User{
		UserName: "admin",
		Password: "123456",
	}
	err = UserDao.Create(user)
	assert.Equal(t, nil, err, "database error")

	user = model.User{
		UserName: "test1",
		Password: "123456",
	}
	err = UserDao.Create(user)
	assert.Equal(t, nil, err, "database error")

	var userInfo *model.User

	userInfo, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "123456", userInfo.Password, "database error")

	err = UserDao.Update(userInfo.ID, map[string]interface{}{
		"Ban": true,
	})
	assert.Equal(t, nil, err, "database error")

	userInfo, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, true, userInfo.Ban, "database error")

	err = UserDao.UpdateByName("test", map[string]interface{}{
		"Ban": false,
	})
	assert.Equal(t, nil, err, "database error")

	userInfo, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, false, userInfo.Ban, "database error")

	var count int64
	count, err = UserDao.UserCount()
	assert.Equal(t, nil, err, "database error")
	assert.Equal(t, int64(3), count, "database error")

	var userList []model.User
	userList, err = UserDao.AllUser()
	assert.Equal(t, nil, err, "database error")
	assert.Equal(t, "test", userList[0].UserName, "database error")

	userList, err = UserDao.GetUsersByNames([]string{"test", "admin"})
	assert.Equal(t, nil, err, "database error")
	assert.Equal(t, "test", userList[0].UserName, "database error")

	err = UserDao.ModifyUserIdentity("haha", 1)
	assert.Equal(t, false, err == nil, "database error")
	err = UserDao.ModifyUserIdentity("test", 6)
	assert.Equal(t, false, err == nil, "database error")
	err = UserDao.ModifyUserIdentity("test", 0)
	assert.Equal(t, nil, err, "database error")
	err = UserDao.ModifyUserIdentity("test", 1)
	assert.Equal(t, nil, err, "database error")
	userInfo, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, true, userInfo.DepartmentSuper, "database error")
	err = UserDao.ModifyUserIdentity("test", 2)
	assert.Equal(t, nil, err, "database error")
	userInfo, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, true, userInfo.EntitySuper, "database error")
	err = UserDao.ModifyUserIdentity("test", 3)
	assert.Equal(t, nil, err, "database error")
	userInfo, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, true, userInfo.SystemSuper, "database error")

	err = UserDao.ModifyUserPassword("haha", "111111")
	assert.Equal(t, false, err == nil, "database error")
	err = UserDao.ModifyUserPassword("test", "654321")
	assert.Equal(t, nil, err, "database error")
	userInfo, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "654321", userInfo.Password, "database error")

	err = UserDao.ModifyUserBanstate("haha", true)
	assert.Equal(t, false, err == nil, "database error")
	err = UserDao.ModifyUserBanstate("test", true)
	assert.Equal(t, nil, err, "database error")
	userInfo, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, true, userInfo.Ban, "database error")
}
