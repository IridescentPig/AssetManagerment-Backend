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

	err = UserDao.Delete([]uint{1})
	assert.Equal(t, nil, err, "database error")
	count, err = UserDao.UserCount()
	assert.Equal(t, nil, err, "database error")
	assert.Equal(t, int64(2), count, "database error")
}

func TestDepartmentEntity(t *testing.T) {
	Init()
	user := model.User{
		UserName: "test",
		Password: "123456",
	}
	department := model.Department{
		Name: "test_department",
	}
	entity := model.Entity{
		Name: "test_entity",
	}
	err := UserDao.Create(user)
	assert.Equal(t, nil, err, "database error")
	err = DepartmentDao.Create(department)
	assert.Equal(t, nil, err, "database error")
	err = EntityDao.Create(entity)
	assert.Equal(t, nil, err, "database error")

	departments, err := DepartmentDao.AllDepartment()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "test_department", departments[0].Name, "database error")
	department = departments[0]

	dc, err := DepartmentDao.DepartmentCount()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, int64(1), dc, "database error")

	ec, err := EntityDao.EntityCount()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, int64(1), ec, "database error")

	entities, err := EntityDao.AllEntity()
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "test_entity", entities[0].Name, "database error")
	entity = entities[0]

	err = UserDao.ModifyUserDepartment("test", department)
	if err != nil {
		log.Fatal(err)
	}

	new_user, err := UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	//log.Print(new_user)
	assert.Equal(t, "test_department", new_user.Department.Name, "database error")

	qd, err := UserDao.GetUserDepartment("test")
	if err != nil {
		log.Fatal(err)
	}
	//log.Print(qd)
	assert.Equal(t, "test_department", qd.Name, "database error")

	err = UserDao.ModifyUserEntity("test", entity)
	if err != nil {
		log.Fatal(err)
	}
	//log.Print(entity)

	new_user, err = UserDao.GetUserByName("test")
	if err != nil {
		log.Fatal(err)
	}
	//log.Print(new_user)
	assert.Equal(t, "test_entity", new_user.Entity.Name, "database error")

	et, err := UserDao.GetUserEntity("test")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "test_entity", et.Name, "database error")

	parent_department := model.Department{
		Name: "parent_department",
	}
	err = DepartmentDao.Create(parent_department)
	assert.Equal(t, nil, err, "database error")

	err = DepartmentDao.ModifyParentDepartment("test_department", "parent_department")
	if err != nil {
		log.Fatal(err)
	}

	/*departments, err = DepartmentDao.AllDepartment()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(departments)*/

	pd, err := DepartmentDao.GetParentDepartment("test_department")
	if err != nil {
		log.Fatal(err)
	}
	//log.Print(pd.Name)
	assert.Equal(t, "parent_department", pd.Name, "database error")

	//new_department, err := DepartmentDao.GetDepartmentByName("test_department")
	//log.Print("parent is: ", new_department.Parent.Name)
}
