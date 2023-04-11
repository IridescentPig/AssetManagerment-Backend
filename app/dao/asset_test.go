package dao

import (
	"asset-management/app/model"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

var database_error string = "database error"

func TestAsset(t *testing.T) {
	Init()
	asset_class := model.AssetClass{
		Name: "test_class",
		Type: 1,
	}

	department := model.Department{
		Name: "test_department",
	}

	err := AssetClassDao.Create(asset_class)
	assert.Equal(t, nil, err, database_error)
	err = DepartmentDao.Create(department)
	assert.Equal(t, nil, err, database_error)

	new_department, err := DepartmentDao.GetDepartmentByID(1)
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
	err = AssetClassDao.ModifyParentAssetClass(int(new_sub_class.ID), int(new_class.ID))
	assert.Equal(t, nil, err, database_error)
	p_class, err := AssetClassDao.GetParentAssetClass(2)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test_class", p_class.Name, database_error)
	s_class, err := AssetClassDao.GetSubAssetClass(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 1, len(s_class), database_error)
	assert.Equal(t, "sub_class", s_class[0].Name, database_error)

	err = AssetClassDao.ModifyAssetClassDepartment(1, int(new_department.ID))
	assert.Equal(t, nil, err, database_error)
	g_dep, err := AssetClassDao.GetAssetClassDepartment(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test_department", g_dep.Name, database_error)

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

	user := model.User{
		UserName: "test1",
		Password: "123456",
	}
	err = UserDao.Create(user)
	assert.Equal(t, nil, err, database_error)

	err = AssetDao.ModifyAssetUser(1, "test1")
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 1, int(new_line.UserID), database_error)
	new_user_s, err := AssetDao.GetAssetUser(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test1", new_user_s.UserName, database_error)

	err = AssetDao.ModifyAssetClass(1, 1)
	assert.Equal(t, nil, err, database_error)
	new_line, err = AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, 1, int(new_line.ClassID), database_error)
	class, err := AssetDao.GetAssetClass(1)
	assert.Equal(t, nil, err, database_error)
	assert.Equal(t, "test_class", class.Name, database_error)

	expire_list := []int{1, 2}
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

}
