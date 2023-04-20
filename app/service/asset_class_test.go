package service

import (
	"asset-management/app/define"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAssetClass(t *testing.T) {
	InitForTest()
	err := EntityService.CreateEntity("test_entity")
	assert.Equal(t, nil, err, "service error")

	err = DepartmentService.CreateDepartment("test_department", 1, 0)
	assert.Equal(t, nil, err, "service error")

	AssetClassService.CreateAssetClass(define.CreateAssetClassReq{
		ClassName: "okok",
		ParentID:  0,
		Type:      1,
	}, 1)

	AssetClassService.ExistsAssetClass(1)
	AssetClassService.ExistsAssetClass(2)
	AssetClassService.GetAssetClassByID(1)
	AssetClassService.GetAssetClassByID(2)

	AssetClassService.CreateAssetClass(define.CreateAssetClassReq{
		ClassName: "sub",
		ParentID:  1,
		Type:      1,
	}, 1)

	AssetClassService.GetSubAssetClass(1, 1)

	parent_id := uint(0)
	AssetClassService.ModifyAssetClassInfo(define.ModifyAssetClassReq{
		ClassName: "new",
		ParentID:  &parent_id,
		Type:      1,
	}, 1)
	parent_id = uint(1)
	AssetClassService.ModifyAssetClassInfo(define.ModifyAssetClassReq{
		ClassName: "new",
		ParentID:  &parent_id,
		Type:      1,
	}, 2)

	AssetClassService.CheckIsAncestor(1, 2)

	AssetClassService.ClassHasAsset(1)

	AssetClassService.DeleteAssetClass(1)

}
