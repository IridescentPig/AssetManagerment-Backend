package dao

import (
	"asset-management/app/model"
	"log"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestUserAsset(t *testing.T) {
	Init()

	user := model.User{
		UserName: "test",
		Password: "123456",
	}
	err := UserDao.Create(user)
	assert.Equal(t, nil, err, "database error")

	p_asset := model.Asset{
		Name:        "parent",
		Price:       decimal.New(1000, 0),
		Description: "test",
		Position:    "OutSpace",
		Number:      10000,
		Type:        1,
		State:       1,
	}
	err = AssetDao.Create(p_asset)
	assert.Equal(t, nil, err, "database error")

	c_asset := model.Asset{
		Name:        "child",
		Price:       decimal.New(1000, 0),
		Description: "test",
		Position:    "OutSpace",
		Number:      10000,
		Type:        1,
		State:       1,
	}
	err = AssetDao.Create(c_asset)
	assert.Equal(t, nil, err, "database error")

	cc_asset := model.Asset{
		Name:        "child_child",
		Price:       decimal.New(1000, 0),
		Description: "test",
		Position:    "OutSpace",
		Number:      10000,
		Type:        1,
		State:       1,
	}
	err = AssetDao.Create(cc_asset)

	AssetDao.ModifyParentAsset(2, 1)
	assert.Equal(t, nil, err, "database error")
	AssetDao.ModifyParentAsset(3, 2)
	assert.Equal(t, nil, err, "database error")

	err = AssetDao.ModifyAssetUser(1, "test")
	assert.Equal(t, nil, err, "database error")
	err = AssetDao.ModifyAssetUser(2, "test")
	assert.Equal(t, nil, err, "database error")
	err = AssetDao.ModifyAssetUser(3, "test")
	assert.Equal(t, nil, err, "database error")

	children, err := AssetDao.GetSubAsset(1)
	assert.Equal(t, nil, err, "database error")
	log.Print(children[0].Name)
	log.Print(children[0].User.UserName)

	parent, err := AssetDao.GetAssetByID(1)
	assert.Equal(t, nil, err, "database error")
	log.Print(parent.Name)
	log.Print(parent.User.UserName)

	direct, err := AssetDao.GetDirectAssetsByUser(1)
	assert.Equal(t, nil, err, "database error")
	log.Print(direct)

	AssetDao.GetDepartmentAssetsByIDs([]uint{0}, 1)
	AssetDao.GetDepartmentAssetsByIDs([]uint{1, 2, 3}, 1)
	AssetDao.GetUserAssetsByIDs([]uint{0}, 1)
	AssetDao.GetUserAssetsByIDs([]uint{1, 2, 3}, 1)
	AssetDao.GetDepartmentAssetsByIDs([]uint{0}, 1)
	AssetDao.GetDepartmentAssetsByIDs([]uint{1, 2, 3}, 1)
	AssetDao.GetDepartmentIdleAssetsByIDs([]uint{0}, 1)
	AssetDao.GetDepartmentIdleAssetsByIDs([]uint{1, 2, 3}, 1)
	AssetDao.ModifyAssetsUserAndState([]uint{0}, 1, 2)
	AssetDao.ModifyAssetsUserAndState([]uint{1, 2, 3}, 1, 3)
	AssetDao.GetUserMaintainAssets(1)
	AssetDao.GetUserMaintainAssets(9)
	AssetDao.ModifyAssetMaintainerAndState([]uint{0}, 1)
	AssetDao.ModifyAssetMaintainerAndState([]uint{1, 2, 3}, 2)

}
