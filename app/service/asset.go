package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/utils"
	"encoding/json"
	"time"

	"github.com/shopspring/decimal"
	"github.com/thoas/go-funk"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type assetService struct{}

var AssetService *assetService

func newAssetService() *assetService {
	return &assetService{}
}

func init() {
	AssetService = newAssetService()
}

func (asset *assetService) TransformAssetBasicInfo(assetList []*model.Asset) []*define.AssetBasicInfo {
	subAssetTreeNodeList := funk.Map(assetList, func(thisAsset *model.Asset) *define.AssetBasicInfo {
		return &define.AssetBasicInfo{
			AssetID:   thisAsset.ID,
			AssetName: thisAsset.Name,
			ParentID:  thisAsset.ParentID,
			User: define.AssetUserBasicInfo{
				UserID:   thisAsset.UserID,
				Username: thisAsset.User.UserName,
			},
			Price:    thisAsset.Price,
			Position: thisAsset.Position,
			Expire:   thisAsset.Expire,
			Class: define.AssetClassBasicInfo{
				ClassID:   thisAsset.ClassID,
				ClassName: thisAsset.Class.Name,
			},
			Number:       thisAsset.Number,
			Type:         thisAsset.Type,
			State:        thisAsset.State,
			Property:     thisAsset.Property,
			NetWorth:     thisAsset.NetWorth,
			DepartmentID: thisAsset.DepartmentID,
			Threshold:    thisAsset.Threshold,
			Warn:         thisAsset.Warn,
		}
	}).([]*define.AssetBasicInfo)

	return subAssetTreeNodeList
}

func (asset *assetService) GetSubAsset(parentID uint, departmentID uint) ([]*define.AssetBasicInfo, error) {
	var subAssetList []*model.Asset
	var err error

	if parentID == 0 {
		subAssetList, err = dao.AssetDao.GetAssetDirectDepartment(departmentID)
	} else {
		subAssetList, err = dao.AssetDao.GetSubAsset(parentID)
	}

	if err != nil {
		return nil, err
	}

	// subAssetTreeNodeList := []*define.AssetInfo{}
	// err = copier.Copy(&subAssetTreeNodeList, subAssetList)

	subAssetTreeNodeList := asset.TransformAssetBasicInfo(subAssetList)

	if err != nil {
		return nil, err
	}

	for _, subNode := range subAssetTreeNodeList {
		subNode.Children, err = asset.GetSubAsset(subNode.AssetID, departmentID)
		if err != nil {
			return nil, err
		}
	}

	return subAssetTreeNodeList, err
}

func (asset *assetService) GetAssetByID(assetID uint) (*model.Asset, error) {
	return dao.AssetDao.GetAssetByID(assetID)
}

func (asset *assetService) ExistAsset(assetID uint) (bool, error) {
	thisAsset, err := dao.AssetDao.GetAssetByID(assetID)
	if err != nil {
		return false, err
	}
	return thisAsset != nil, nil
}

func (asset *assetService) CheckAssetInDepartment(assetID uint, departmentID uint) (bool, error) {
	thisAsset, err := dao.AssetDao.GetAssetByID(assetID)
	if err != nil || thisAsset == nil {
		return false, nil
	}
	return thisAsset.DepartmentID == departmentID, nil
}

func (asset *assetService) CheckIsAncestor(srcID uint, targetID uint) (bool, error) {
	targetAsset, err := dao.AssetDao.GetAssetByID(targetID)
	if err != nil {
		return true, err
	}

	flag := false
	for {
		if targetAsset == nil {
			break
		}
		if targetAsset.ID == srcID {
			flag = true
			break
		}
		targetAsset, err = dao.AssetDao.GetAssetByID(targetAsset.ParentID)
		if err != nil {
			return true, err
		}
	}

	return flag, nil
}

func (asset *assetService) ModifyAssetInfo(id uint, req define.ModifyAssetInfoReq) error {
	err := dao.AssetDao.UpdateByStruct(id, model.Asset{
		Name:        req.AssetName,
		Price:       req.Price,
		Description: req.Description,
		Position:    req.Position,
		ClassID:     req.ClassID,
		Type:        req.Type,
		Number:      req.Number,
		Expire:      req.Expire,
		ImgList:     req.ImgList,
		Threshold:   req.Threshold,
	})
	if err != nil {
		return err
	}
	if req.ParentID != nil {
		if *req.ParentID != 0 {
			err = dao.AssetDao.Update(id, map[string]interface{}{
				"parent_id": *req.ParentID,
			})
		} else {
			err = dao.AssetDao.Update(id, map[string]interface{}{
				"parent_id": gorm.Expr("NULL"),
			})
		}
	}
	return err
}

/*
Update net worth and warning tag
*/
func (asset *assetService) UpdateNetWorth(assetID uint) error {
	thisAsset, err := dao.AssetDao.GetAssetByID(assetID)
	if err != nil {
		return err
	} else if thisAsset == nil {
		return nil
	}

	if thisAsset.Expire == 0 || thisAsset.State >= 3 {
		return nil
	}

	price := thisAsset.Price
	expire := thisAsset.Expire
	interval := utils.GetDiffDays(time.Time(*thisAsset.CreatedAt), time.Now())

	if interval >= int(thisAsset.Expire) {
		err = dao.AssetDao.Update(assetID, map[string]interface{}{
			"net_worth": decimal.Zero,
			"state":     3,
		})
		if err != nil {
			return err
		}

		subAssets, err := dao.AssetDao.GetSubAsset(assetID)
		if err != nil {
			return err
		}

		if subAssets != nil {
			subIds := funk.Map(subAssets, func(currentAsset *model.Asset) uint {
				return currentAsset.ID
			}).([]uint)

			err = dao.AssetDao.AllUpdate(subIds, map[string]interface{}{
				"parent_id": gorm.Expr("NULL"),
			})

			if err != nil {
				return err
			}
		}
	} else {
		rate := 1.0 - float64(interval)/float64(expire)
		isWarn := (int(expire) - interval) <= int(thisAsset.Threshold)
		err = dao.AssetDao.UpdateByStruct(assetID, model.Asset{
			NetWorth: price.Mul(decimal.NewFromFloat(rate)),
			Warn:     isWarn,
		})
	}

	return err
}

func (asset *assetService) CreateAsset(req *define.CreateAssetReq, departmentID uint, parentID uint, userID uint) error {
	thisID, err := dao.AssetDao.CreateAndGetID(model.Asset{
		Name:         req.AssetName,
		Price:        req.Price,
		Description:  req.Description,
		Position:     req.Position,
		ClassID:      req.ClassID,
		Number:       req.Number,
		Type:         req.Type,
		DepartmentID: departmentID,
		UserID:       userID,
		ParentID:     parentID,
		Property:     datatypes.JSON([]byte(`{}`)),
		Expire:       req.Expire,
		NetWorth:     req.Price,
		ImgList:      req.ImgList,
		Threshold:    req.Threshold,
	})
	if err != nil {
		return err
	}
	for _, child := range req.Children {
		err = asset.CreateAsset(child, departmentID, thisID, userID)
		if err != nil {
			return err
		}
	}
	return nil
}

func (asset *assetService) ExpireAssets(assetIDs []uint) error {
	err := dao.AssetDao.ExpireAsset(assetIDs)
	return err
}

func (asset *assetService) TransferAssets(assetIDs []uint, userID uint, departmentID uint, oldDepartmentID uint) error {
	if departmentID != oldDepartmentID {
		subAssets, err := dao.AssetDao.GetSubAssetsByParents(assetIDs)
		if err != nil {
			return err
		}
		subAssetIDs := []uint{}

		for _, asset := range subAssets {
			subAssetIDs = append(subAssetIDs, asset.ID)
		}

		err = dao.AssetDao.AllUpdate(subAssetIDs, map[string]interface{}{
			"parent_id": gorm.Expr("NULL"),
		})
		if err != nil {
			return err
		}
		err = dao.AssetDao.AllUpdate(assetIDs, map[string]interface{}{
			"user_id":       userID,
			"parent_id":     gorm.Expr("NULL"),
			"department_id": departmentID,
		})
		return err
	} else {
		err := dao.AssetDao.AllUpdate(assetIDs, map[string]interface{}{
			"user_id": userID,
		})
		return err
	}
}

func (asset *assetService) GetAssetByUser(userID uint) (assets []*define.AssetBasicInfo, err error) {
	assetList, err := dao.AssetDao.GetDirectAssetsByUser(userID)
	if err != nil {
		return
	}

	assets = asset.TransformAssetBasicInfo(assetList)
	for _, childAsset := range assets {
		childAsset.Children, err = asset.GetSubAsset(childAsset.AssetID, childAsset.DepartmentID)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	return
}

func (asset *assetService) GetDepartmentAssetsByIDs(ids []uint, departmentID uint) ([]*model.Asset, error) {
	assetList, err := dao.AssetDao.GetDepartmentAssetsByIDs(ids, departmentID)
	if err != nil {
		return nil, err
	}
	return assetList, nil
}

func (asset *assetService) GetUserAssetsByIDs(ids []uint, userID uint) ([]*model.Asset, error) {
	assetList, err := dao.AssetDao.GetUserAssetsByIDs(ids, userID)
	if err != nil {
		return nil, err
	}
	return assetList, nil
}

func (asset *assetService) GetDepartmentIdleAssets(ids []uint, departmentID uint) ([]*model.Asset, error) {
	assetList, err := dao.AssetDao.GetDepartmentIdleAssetsByIDs(ids, departmentID)
	if err != nil {
		return nil, err
	}
	return assetList, nil
}

func (asset *assetService) AcquireAssets(ids []uint, userID uint) error {
	err := dao.AssetDao.ModifyAssetsUserAndState(ids, userID, 1)
	return err
}

func (asset *assetService) CancelAssets(ids []uint, userID uint) error {
	err := dao.AssetDao.ModifyAssetsUserAndState(ids, userID, 0)
	return err
}

func (asset *assetService) GetUserMaintainAssets(userID uint) ([]*model.Asset, error) {
	assetList, err := dao.AssetDao.GetUserMaintainAssets(userID)
	return assetList, err
}

func (asset *assetService) ModifyAssetMaintainerAndState(assetIDs []uint, maintainerID uint) error {
	err := dao.AssetDao.ModifyAssetMaintainerAndState(assetIDs, maintainerID)
	return err
}

func (asset *assetService) ExistsProperty(assetID uint, key string) (bool, error) {
	return dao.AssetDao.CheckAssetPropertyExist(assetID, key)
}

func (asset *assetService) SetProperty(assetID uint, key string, value string) error {
	return dao.AssetDao.SetAssetProperty(assetID, key, value)
}

func (asset *assetService) DeleteProperty(assetID uint, key string) error {
	thisAsset, err := dao.AssetDao.GetAssetProperty(assetID)
	if err != nil {
		return err
	}

	var data map[string]interface{}
	err = json.Unmarshal(thisAsset.Property, &data)
	if err != nil {
		return err
	}

	delete(data, key)

	jsonData, err := json.Marshal(data)
	if err != nil {
		return err
	}

	err = dao.AssetDao.Update(assetID, map[string]interface{}{
		"property": jsonData,
	})

	return err
}

func (asset *assetService) GetAssetHistory(assetID uint) ([]*model.Task, error) {
	taskList, err := dao.AssetDao.GetAssetTask(assetID)
	if err != nil {
		return nil, err
	}

	var approvedTaskList []*model.Task

	for _, task := range taskList {
		if task.State == 1 {
			approvedTaskList = append(approvedTaskList, task)
		}
	}

	return approvedTaskList, nil
}

func (asset *assetService) SearchDepartmentAssets(departmentID uint, req *define.SearchAssetReq) ([]*model.Asset, error) {
	if req.Name != "" {
		req.Name = "%" + req.Name + "%"
	}

	if req.Description != "" {
		req.Description = "%" + req.Description + "%"
	}

	return dao.AssetDao.SearchDepartmentAsset(departmentID, req)
}

func (asset *assetService) GetDepartmentAssetCount(departmentID uint) (int64, error) {
	return dao.AssetDao.GetDepartmentAssetCount(departmentID)
}

func (asset *assetService) GetDepartmentAssetInWarn(departmentID uint) ([]*model.Asset, error) {
	return dao.AssetDao.GetDepartmentWarnAsset(departmentID)
}
