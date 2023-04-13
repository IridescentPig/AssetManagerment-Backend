package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"

	"github.com/jinzhu/copier"
	"gorm.io/gorm"
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
	err := dao.AssetClassDao.UpdateByStruct(id, model.AssetClass{
		Name: req.ClassName,
		Type: req.Type,
	})
	if err != nil {
		return err
	}
	if req.ParentID != nil {
		if *req.ParentID != 0 {
			err = dao.AssetClassDao.Update(id, map[string]interface{}{
				"parent_id": *req.ParentID,
			})
		} else {
			err = dao.AssetClassDao.Update(id, map[string]interface{}{
				"parent_id": gorm.Expr("NULL"),
			})
		}
	}
	return err
}

func (assetClass *assetClassService) CheckIsAncestor(srcClassID uint, targetClassID uint) (bool, error) {
	targetClass, err := dao.AssetClassDao.GetAssetClassByID(targetClassID)
	if err != nil {
		return false, nil
	}

	flag := false
	for {
		if targetClass == nil {
			break
		}
		if targetClass.ID == srcClassID {
			flag = true
			break
		}
		targetClass, err = dao.AssetClassDao.GetAssetClassByID(targetClass.ParentID)
		if err != nil {
			return true, err
		}
	}
	return flag, nil
}

func (assetClass *assetClassService) ClassHasAsset(classID uint) (bool, error) {
	assetList, err := dao.AssetDao.GetAssetListByClassID(classID)
	if err != nil {
		return true, err
	}
	return len(assetList) != 0, nil
}

func (assetClass *assetClassService) DeleteAssetClass(classID uint) error {
	return dao.AssetClassDao.Delete([]uint{classID})
}

func (assetClass *assetClassService) GetSubClass(parentID uint, departmentID uint) ([]*define.AssetClassTreeNode, error) {
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

	return subTreeNodeList, nil
}
