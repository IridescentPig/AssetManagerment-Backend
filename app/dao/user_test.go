package dao

import (
	"asset-management/app/model"
	"log"
	"testing"

	"github.com/shopspring/decimal"
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

	var userList []*model.User
	userList, err = UserDao.AllUser()
	assert.Equal(t, nil, err, "database error")
	assert.Equal(t, "test", userList[0].UserName, "database error")

	var users []model.User
	users, err = UserDao.GetUsersByNames([]string{"test", "admin"})
	assert.Equal(t, nil, err, "database error")
	assert.Equal(t, "test", users[0].UserName, "database error")

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
	// log.Println(new_user.Department)
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
	// log.Print(new_user.Entity)
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

	sd, err := DepartmentDao.GetSubDepartment("parent_department")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 1, len(sd), "database error")
	assert.Equal(t, "test_department", sd[0].Name, "database error")

	err = DepartmentDao.ModifyDepartmentEntity("parent_department", "test_entity")
	if err != nil {
		log.Fatal(err)
	}

	/*departments, err = DepartmentDao.AllDepartment()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(departments)

	entities, err = EntityDao.AllEntity()
	if err != nil {
		log.Fatal(err)
	}
	log.Print(entities)*/

	se, err := DepartmentDao.GetDepartmentEntity("parent_department")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, "test_entity", se.Name, "database error")

	du, err := DepartmentDao.GetDepartmentDirectUser("parent_department")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 0, len(du), "database error")

	au, err := DepartmentDao.GetDepartmentAllUser("parent_department")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 1, len(au), "database error")
	assert.Equal(t, "test", au[0].UserName, "database error")

	au, err = EntityDao.GetEntityAllUser(1)
	if err != nil {
		log.Fatal(err)
	}
	//log.Print("au: ", au)
	//log.Print(au[0].UserName)
	assert.Equal(t, 1, len(au), "database error")
	assert.Equal(t, "test", au[0].UserName, "database error")

	ad, err := EntityDao.GetEntityAllDepartment(1)
	if err != nil {
		log.Fatal(err)
	}
	//log.Print("ad: ", ad)
	//log.Print(ad[0].Name)
	assert.Equal(t, 1, len(ad), "database error")
	assert.Equal(t, "parent_department", ad[0].Name, "database error")

	sd, err = EntityDao.GetEntitySubDepartment("test_entity")
	if err != nil {
		log.Fatal(err)
	}
	//log.Print("sd: ", sd)
	//log.Print(sd[0].Name)
	assert.Equal(t, 1, len(sd), "database error")
	assert.Equal(t, "parent_department", sd[0].Name, "database error")

	DepartmentDao.GetDepartmentSub("test_department", 1, 1)
	DepartmentDao.GetSubDepartmentByID(1)
	DepartmentDao.GetDepartmentDirectUserByID(1)
	DepartmentDao.GetDepartmentDirectUser("test_department")
	DepartmentDao.GetDepartmentAllUserByID(1)
	DepartmentDao.GetDepartmentManager(1)
	DepartmentDao.Update(1, map[string]interface{}{
		"password": "123098439",
	})
	DepartmentDao.GetDepartmentsByNames([]string{"test_department"})
	DepartmentDao.Delete([]uint{1})

	EntityDao.GetEntitysByNames([]string{"test_entity"})
	EntityDao.GetEntityByID(1)
	EntityDao.GetEntityManager(1)
	EntityDao.GetEntitySubDepartmentByID(1)
	EntityDao.GetEntitySubDepartmentByID(9)
	EntityDao.Update(1, map[string]interface{}{
		"name": "askdjhjka",
	})
	EntityDao.Delete([]uint{1})

	UserDao.GetLimitUser(1, 10)
	UserDao.GetUserByID(1)
	UserDao.ModifyUserEntityByID(1, 1)
	UserDao.ModifyUserDepartmentByID(1, 1)

}

var database_error string = "database error"

func TestAsset(t *testing.T) {
	Init()

	asset_class := model.AssetClass{
		Name: "test_class",
		Type: 1,
	}

	department := model.Department{
		Name: "test_asset_department",
	}

	err := AssetClassDao.Create(asset_class)
	assert.Equal(t, nil, err, database_error)
	err = DepartmentDao.Create(department)
	assert.Equal(t, nil, err, database_error)

	new_department, err := DepartmentDao.GetDepartmentByID(3)
	assert.Equal(t, nil, err, database_error)

	new_class, err := AssetClassDao.GetAssetClassByID(1)
	assert.Equal(t, nil, err, database_error)

	sub_class := model.AssetClass{
		Name: "sub_class",
		Type: 2,
	}
	err = AssetClassDao.Create(sub_class)
	assert.Equal(t, nil, err, database_error)
	new_sub_class, err := AssetClassDao.GetAssetClassByID(2)
	assert.Equal(t, nil, err, database_error)
	err = AssetClassDao.ModifyParentAssetClass(new_sub_class.ID, new_class.ID)
	assert.Equal(t, nil, err, database_error)
	p_class, err := AssetClassDao.GetParentAssetClass(2)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test_class", p_class.Name, database_error)
	s_class, err := AssetClassDao.GetSubAssetClass(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 1, len(s_class), database_error)
	assert.Equal(t, "sub_class", s_class[0].Name, database_error)

	err = AssetClassDao.ModifyAssetClassDepartment(1, new_department.ID)
	assert.Equal(t, nil, err, database_error)
	g_dep, err := AssetClassDao.GetAssetClassDepartment(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test_asset_department", g_dep.Name, database_error)

	line_asset := model.Asset{
		Name:        "test_asset_line",
		Price:       decimal.New(100, 0),
		Description: "test",
		Position:    "OffSpace",
		Number:      1,
		Type:        1,
	}

	num_asset := model.Asset{
		Name:        "test_asset_num",
		Price:       decimal.New(1000, 0),
		Description: "test",
		Position:    "OutSpace",
		Number:      10000,
		Type:        1,
	}

	err = AssetDao.Create(line_asset)
	assert.Equal(t, nil, err, database_error)
	err = AssetDao.Create(num_asset)
	assert.Equal(t, nil, err, database_error)

	list, err := AssetDao.AllAsset()
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 2, len(list), database_error)

	new_line, err := AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test_asset_line", new_line.Name, database_error)

	new_num, err := AssetDao.GetAssetByName("test_asset_num")
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 1, len(new_num), database_error)
	assert.Equal(t, "test_asset_num", new_num[0].Name, database_error)

	count, err := AssetDao.AssetCount()
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, int64(2), count, database_error)

	err = AssetDao.ModifyAssetPrice(1, decimal.NewFromFloat(233))
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, decimal.NewFromFloat(233), new_line.Price, database_error)
	err = AssetDao.ModifyAssetDescription(1, "new")
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "new", new_line.Description, database_error)
	err = AssetDao.ModifyAssetPosition(1, "M78")
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "M78", new_line.Position, database_error)
	err = AssetDao.ModifyAssetNum(2, 233)
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(2)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 233, new_line.Number, database_error)

	user := model.User{
		UserName: "test_asset_user",
		Password: "123456",
	}
	err = UserDao.Create(user)
	assert.Equal(t, nil, err, database_error)

	err = AssetDao.ModifyAssetUser(1, "test_asset_user")
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 1, int(new_line.UserID), database_error)
	new_user_s, err := AssetDao.GetAssetUser(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test_asset_user", new_user_s.UserName, database_error)

	err = AssetDao.ModifyParentAsset(2, 1)
	assert.Equal(t, nil, err, database_error)
	new_num_s, err := AssetDao.GetAssetByID(2)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, int(new_num_s.ParentID), 1, database_error)
	new_line, err = AssetDao.GetParentAsset(2)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, new_line.Name, "test_asset_line", database_error)
	sub_s, err := AssetDao.GetSubAsset(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 1, len(sub_s), database_error)
	assert.Equal(t, sub_s[0].Name, "test_asset_num", database_error)

	err = AssetDao.ModifyAssetClass(1, 1)
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 1, int(new_line.ClassID), database_error)
	class, err := AssetDao.GetAssetClass(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test_class", class.Name, database_error)

	expire_list := []uint{1, 2}
	err = AssetDao.ExpireAsset(expire_list)
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	new_num_s, err = AssetDao.GetAssetByID(2)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, true, new_line.Expire, database_error)
	assert.Equal(t, decimal.New(0, 0), new_line.Price, database_error)
	assert.Equal(t, true, new_num_s.Expire, database_error)
	assert.Equal(t, decimal.New(0, 0), new_num_s.Price, database_error)

	AssetClassDao.Update(1, map[string]interface{}{
		"name": "asdkfjhjk",
	})
	AssetClassDao.UpdateByStruct(1, model.AssetClass{
		Name: "test_class",
		Type: 1,
	})
	AssetClassDao.AllUpdate([]uint{1}, map[string]interface{}{
		"name": "asdkfjhjk",
	})
	AssetClassDao.Delete([]uint{1})
	AssetClassDao.GetDepartmentDirectClass(1)

	line_asset = model.Asset{
		Name:        "test_asset_line",
		Price:       decimal.New(100, 0),
		Description: "test",
		Position:    "OffSpace",
		Number:      1,
		Type:        1,
	}

	AssetDao.CreateAndGetID(line_asset)
	AssetDao.UpdateByStruct(1, line_asset)
	AssetDao.Delete([]uint{1})
	AssetDao.GetAssetListByClassID(1)
	AssetDao.GetSubAssetsByParents([]uint{1})
}
