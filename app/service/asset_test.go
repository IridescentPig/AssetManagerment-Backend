package service

import (
	"asset-management/app/define"
	"testing"

	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
)

func TestAsset(t *testing.T) {
	InitForTest()

	_ = EntityService.CreateEntity("test_entity")
	//assert.Equal(t, nil, err, "service error")

	err := DepartmentService.CreateDepartment("test_department", 1, 0)
	assert.Equal(t, nil, err, "service error")

	err = DepartmentService.CreateDepartmentUser(define.CreateDepartmentUserReq{
		UserName:        "test_user",
		Password:        "123456",
		DepartmentSuper: false,
	}, 1, 1)
	assert.Equal(t, nil, err, "service error")

	AssetClassService.CreateAssetClass(define.CreateAssetClassReq{
		ClassName: "okok",
		ParentID:  0,
		Type:      1,
	}, 1)

	CreateAsset2 := define.CreateAssetReq{
		AssetName:   "sub",
		Price:       decimal.New(123, 0),
		Description: "oo",
		Position:    "bb",
		ClassID:     1,
		Number:      1,
		Type:        1,
	}
	CreateAsset := define.CreateAssetReq{
		AssetName:   "name",
		Price:       decimal.New(123, 0),
		Description: "oo",
		Position:    "bb",
		ClassID:     1,
		Number:      1,
		Type:        1,
		Children:    []*define.CreateAssetReq{&CreateAsset2},
	}

	AssetService.CreateAsset(&CreateAsset, 1, 0, 1)

	parent_id := uint(0)
	ModifyAsset := define.ModifyAssetInfoReq{
		AssetName:   "name",
		ParentID:    &parent_id,
		Price:       decimal.New(123, 0),
		Description: "oo",
		Position:    "bb",
		ClassID:     1,
		Number:      1,
		Type:        1,
	}

	AssetService.ModifyAssetInfo(1, ModifyAsset)

	AssetService.CreateAsset(&CreateAsset, 1, 1, 1)

	AssetService.GetSubAsset(1, 1)
	AssetService.GetSubAsset(0, 1)

	AssetService.GetAssetByID(1)
	AssetService.ExistAsset(1)
	AssetService.CheckAssetInDepartment(1, 1)
	AssetService.CheckIsAncestor(1, 2)
	assets := []uint{1}
	AssetService.ExpireAssets(assets)
	assets = []uint{2}
	AssetService.ExpireAssets(assets)
	AssetService.TransferAssets([]uint{1}, 1, 1, 1)
	AssetService.TransferAssets([]uint{1, 2, 3}, 1, 1, 2)
	AssetService.TransferAssets([]uint{1, 2, 3}, 1, 2, 1)

	AssetService.GetAssetByUser(1)
	AssetService.GetAssetByUser(2)
	AssetService.GetAssetByUser(90)
	AssetService.GetDepartmentAssetsByIDs([]uint{0}, 1)
	AssetService.GetDepartmentAssetsByIDs([]uint{1, 2, 3}, 1)
	AssetService.GetUserAssetsByIDs([]uint{0}, 1)
	AssetService.GetUserAssetsByIDs([]uint{1, 2, 3}, 1)
	AssetService.GetDepartmentIdleAssets([]uint{0}, 1)
	AssetService.GetDepartmentIdleAssets([]uint{1, 2, 3}, 1)
	AssetService.AcquireAssets([]uint{0}, 1)
	AssetService.AcquireAssets([]uint{1, 2, 3}, 1)
	AssetService.CancelAssets([]uint{0}, 1)
	AssetService.CancelAssets([]uint{1, 2, 3}, 1)
	AssetService.GetUserMaintainAssets(1)
	AssetService.GetUserMaintainAssets(9)
	AssetService.ModifyAssetMaintainerAndState([]uint{0}, 1)
	AssetService.ModifyAssetMaintainerAndState([]uint{1, 2, 3}, 1)

}
