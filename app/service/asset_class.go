package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"

	"github.com/jinzhu/copier"
)

type assetClassService struct{}

var AssetClassService *assetClassService

func newAssetClassService() *assetClassService {
	return &assetClassService{}
}

func init() {
	AssetClassService = newAssetClassService()
}

func (assetClass *assetClassService) GetAssetClassByID(id uint) (*model.AssetClass, error) {
	thisAssetClass, err := dao.AssetClassDao.GetAssetClassByID(id)
	return thisAssetClass, err
}

func (assetClass *assetClassService) ExistsAssetClass(id uint) (bool, error) {
	thisAssetClass, err := assetClass.GetAssetClassByID(id)
	if err != nil {
		return false, err
	}
	return thisAssetClass != nil, err
}

func (assetClass *assetClassService) CreateAssetClass(req define.CreateAssetClassReq, departmentID uint) error {
	return dao.AssetClassDao.Create(model.AssetClass{
		Name:         req.ClassName,
		ParentID:     req.ParentID,
		DepartmentID: departmentID,
		Type:         req.Type,
	})
}

func (assetClass *assetClassService) GetSubAssetClass(parentID uint, departmentID uint) ([]*define.AssetClassTreeNode, error) {
	var subClassList []*model.AssetClass
	var err error
	if parentID == 0 {
		subClassList, err = dao.AssetClassDao.GetDepartmentDirectClass(departmentID)
	} else {
		subClassList, err = dao.AssetClassDao.GetSubAssetClass(parentID)
	}

	if err != nil {
		return nil, err
	}
	subTreeNodeList := []*define.AssetClassTreeNode{}
	err = copier.Copy(&subTreeNodeList, subClassList)
	if err != nil {
		return nil, err
	}

	for _, subNode := range subTreeNodeList {
		subNode.Children, err = assetClass.GetSubAssetClass(subNode.ClassID, departmentID)
		if err != nil {
			return nil, err
		}
	}

	return subTreeNodeList, nil
}

func (assetClass *assetClassService) ModifyAssetClassInfo(req define.ModifyAssetClassReq, id uint) error {
	return dao.AssetClassDao.UpdateByStruct(id, model.AssetClass{
		Name:     req.ClassName,
		ParentID: req.ParentID,
		Type:     req.Type,
	})
}
