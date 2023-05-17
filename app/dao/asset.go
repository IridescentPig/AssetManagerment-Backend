package dao

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"asset-management/utils"
	"errors"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type assetDao struct {
}

var AssetDao *assetDao

func newAssetDao() *assetDao {
	return &assetDao{}
}

func init() {
	AssetDao = newAssetDao()
}

func (asset *assetDao) Create(newAsset model.Asset) error {
	result := db.Model(&model.Asset{}).Create(&newAsset)
	return utils.DBError(result)
}

func (asset *assetDao) CreateAndGetID(newAsset model.Asset) (uint, error) {
	result := db.Model(&model.Asset{}).Create(&newAsset)
	return newAsset.ID, utils.DBError(result)
}

func (asset *assetDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.Asset{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (asset *assetDao) UpdateByStruct(id uint, data model.Asset) error {
	result := db.Model(model.Asset{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (asset *assetDao) AllUpdate(ids []uint, data map[string]interface{}) error {
	result := db.Model(&model.Asset{}).Where("id IN (?)", ids).Updates(data)
	return utils.DBError(result)
}

func (asset *assetDao) Delete(id []uint) error {
	result := db.Model(&model.Asset{}).Where("id in (?)", id).Delete(&model.Asset{})
	return utils.DBError(result)
}

func (asset *assetDao) SaveAsset(thisAsset *model.Asset) error {
	result := db.Save(thisAsset)
	return utils.DBError(result)
}

// func (asset *assetDao) AllAsset() (list []model.Asset, err error) {
// 	result := db.Model(&model.Asset{}).Find(&list)
// 	for _, asset := range list {
// 		user := &model.User{}
// 		err = db.Model(&asset).Association("User").Find(&user)
// 		if err != nil {
// 			return
// 		}
// 		asset.UserID = user.ID
// 		asset_class := &model.AssetClass{}
// 		err = db.Model(&asset).Association("Class").Find(&asset_class)
// 		if err != nil {
// 			return
// 		}
// 		asset.ClassID = asset_class.ID
// 	}
// 	err = utils.DBError(result)
// 	return
// }

/*
Note: This function doesn't preload any association function
*/
func (asset *assetDao) GetAllAssets() (assetList []*model.Asset, err error) {
	result := db.Model(&model.Asset{}).Find(&assetList)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}

	return
}

func (asset *assetDao) GetAssetByName(name string) (list []model.Asset, err error) {
	result := db.Model(&model.Asset{}).Where("name = ?", name).Find(&list)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	for _, asset := range list {
		user := &model.User{}
		err = db.Model(&asset).Association("User").Find(&user)
		if err != nil {
			return
		}
		asset.UserID = user.ID
		asset_class := &model.AssetClass{}
		err = db.Model(&asset).Association("Class").Find(&asset_class)
		if err != nil {
			return
		}
		asset.ClassID = asset_class.ID
	}
	err = utils.DBError(result)
	return
}

func (asset *assetDao) GetAssetByID(id uint) (*model.Asset, error) {
	ret := &model.Asset{}
	result := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("ID = ?", id).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	// user := &model.User{}
	// err := db.Model(&ret).Association("User").Find(&user)
	// if err != nil {
	// 	return nil, err
	// }
	// ret.User = *user
	// asset_class := &model.AssetClass{}
	// err = db.Model(&ret).Association("Class").Find(&asset_class)
	// if err != nil {
	// 	return nil, err
	// }
	// ret.Class = *asset_class
	return ret, utils.DBError(result)
}

func (asset *assetDao) AssetCount() (count int64, err error) {
	result := db.Model(&model.Asset{}).Count(&count)
	err = utils.DBError(result)
	return
}

var asset_not_exist string = "asset doesn't exist"

func (asset *assetDao) ModifyAssetPrice(id uint, price decimal.Decimal) error {
	thisAsset, err := asset.GetAssetByID(id)
	if err != nil {
		return err
	}
	if thisAsset == nil {
		return errors.New(asset_not_exist)
	}
	err = asset.Update(thisAsset.ID, map[string]interface{}{
		"price": price,
	})
	return err
}

func (asset *assetDao) ModifyAssetDescription(id uint, description string) error {
	thisAsset, err := asset.GetAssetByID(id)
	if err != nil {
		return err
	}
	if thisAsset == nil {
		return errors.New(asset_not_exist)
	}
	err = asset.Update(thisAsset.ID, map[string]interface{}{
		"description": description,
	})
	return err
}

func (asset *assetDao) ModifyAssetPosition(id uint, position string) error {
	thisAsset, err := asset.GetAssetByID(id)
	if err != nil {
		return err
	}
	if thisAsset == nil {
		return errors.New(asset_not_exist)
	}
	err = asset.Update(thisAsset.ID, map[string]interface{}{
		"position": position,
	})
	return err
}

func (asset *assetDao) ModifyAssetNum(id uint, num int) error {
	thisAsset, err := asset.GetAssetByID(id)
	if err != nil {
		return err
	}
	if thisAsset == nil {
		return errors.New(asset_not_exist)
	}
	err = asset.Update(thisAsset.ID, map[string]interface{}{
		"Number": num,
	})
	return err
}

func (asset *assetDao) ModifyAssetState(id uint, state uint) error {
	thisAsset, err := asset.GetAssetByID(id)
	if err != nil {
		return err
	}
	if thisAsset == nil {
		return errors.New(asset_not_exist)
	}
	err = asset.Update(thisAsset.ID, map[string]interface{}{
		"State": state,
	})
	return err
}

func (asset *assetDao) ExpireAsset(ids []uint) error {
	return asset.AllUpdate(ids, map[string]interface{}{
		"state": 3,
		"price": decimal.NewFromFloat(0),
	})
}

// asset and asset
func (asset *assetDao) GetSubAsset(id uint) (assets []*model.Asset, err error) {
	err = utils.DBError(db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("parent_id = ?", id).Find(&assets))
	return
}

func (asset *assetDao) GetAssetDirectDepartment(departmentID uint) (assets []*model.Asset, err error) {
	err = utils.DBError(db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("department_id = ? and parent_id IS NULL", departmentID).Find(&assets))
	return
}

func (asset *assetDao) GetParentAsset(id uint) (ParentAsset *model.Asset, err error) {
	query_asset, err := asset.GetAssetByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&query_asset).Where("id = ?", query_asset.ParentID).Find(&ParentAsset))
	return
}

func (asset *assetDao) ModifyParentAsset(ChildID uint, ParentID uint) error {
	child_asset, err := asset.GetAssetByID(ChildID)
	if err != nil {
		return err
	}
	parent_asset, err := asset.GetAssetByID(ParentID)
	if err != nil {
		return err
	}
	child_asset.ParentID = parent_asset.ID
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&child_asset))
}

// asset and user
func (asset *assetDao) GetAssetUser(id uint) (user model.User, err error) {
	query_asset, err := asset.GetAssetByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&model.User{}).Where("id = ?", query_asset.UserID).Find(&user))
	return
}

func (asset *assetDao) ModifyAssetUser(AssetID uint, Username string) error {
	query_asset, err := asset.GetAssetByID(AssetID)
	if err != nil {
		return err
	}
	target_user, err := UserDao.GetUserByName(Username)
	if err != nil {
		return err
	}
	query_asset.UserID = target_user.ID
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&query_asset))
}

func (asset *assetDao) GetDirectAssetsByUser(userID uint) (assets []*model.Asset, err error) {
	/*result := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
	Preload("Department").Preload("Class").Preload("Maintainer").
	Where("user_id = ? and parent_id = 0", userID).Find(&assets)*/
	//db = db.Debug()
	sub_query := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("user_id = ?", userID).Select("id")
	var all_assets []uint
	sub_query.Find(&all_assets)
	//log.Print("all: ", len(all_assets), all_assets[0], all_assets[1])

	result := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("user_id = ?", userID).Where("parent_id not in (?) or parent_id is null", all_assets).Find(&assets)
	//log.Print("fliter: ", len(assets) /*assets[0].ParentID, assets[1].ParentID*/)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

// asset and asset_class
func (asset *assetDao) GetAssetClass(id uint) (class model.AssetClass, err error) {
	query_asset, err := asset.GetAssetByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&model.AssetClass{}).Where("id = ?", query_asset.ClassID).Find(&class))
	return
}

var type_not_match = "type not match"

func (asset *assetDao) ModifyAssetClass(AssetID uint, ClassID uint) error {
	query_asset, err := asset.GetAssetByID(AssetID)
	if err != nil {
		return err
	}
	target_class, err := AssetClassDao.GetAssetClassByID(ClassID)
	if err != nil {
		return err
	}
	if query_asset.Type != target_class.Type {
		return errors.New(type_not_match)
	}
	query_asset.ClassID = target_class.ID
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&query_asset))
}

func (asset *assetDao) GetAssetListByClassID(assetClassID uint) ([]*model.Asset, error) {
	var assetList []*model.Asset
	err := utils.DBError(db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("class_id = ?", assetClassID).Find(&assetList))
	return assetList, err
}

func (asset *assetDao) GetSubAssetsByParents(ids []uint) (assets []*model.Asset, err error) {
	result := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("parent_id IN (?)", ids).Find(&assets)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

func (asset *assetDao) GetDepartmentAssetsByIDs(ids []uint, departmentID uint) (assets []*model.Asset, err error) {
	result := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("id IN (?) and department_id = ? and state <= ?", ids, departmentID, 2).Find(&assets)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

func (asset *assetDao) GetUserAssetsByIDs(ids []uint, userID uint) (assets []*model.Asset, err error) {
	result := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("id IN (?) and user_id = ? and state = ?", ids, userID, 1).Find(&assets)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

func (asset *assetDao) GetDepartmentIdleAssetsByIDs(ids []uint, departmentID uint) (assets []*model.Asset, err error) {
	result := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("id IN (?) and department_id = ? and state = ?", ids, departmentID, 0).Find(&assets)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

func (asset *assetDao) ModifyAssetsUserAndState(ids []uint, userID uint, state uint) error {
	result := db.Model(&model.Asset{}).Where("id IN (?)", ids).Updates(map[string]interface{}{
		"user_id": userID,
		"state":   state,
	})
	return utils.DBError(result)
}

func (asset *assetDao) GetUserMaintainAssets(userID uint) (assetList []*model.Asset, err error) {
	result := db.Model(&model.Asset{}).Preload("Parent").Preload("User").
		Preload("Department").Preload("Class").Preload("Maintainer").
		Where("maintainer_id = ?", userID).Find(&assetList)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

func (asset *assetDao) ModifyAssetMaintainerAndState(assetIDs []uint, maintainerID uint) error {
	var result *gorm.DB
	if maintainerID == 0 {
		result = db.Model(&model.Asset{}).Where("id IN (?)", assetIDs).Updates(map[string]interface{}{
			"maintainer_id": gorm.Expr("NULL"),
			"state":         1,
		})
	} else {
		result = db.Model(&model.Asset{}).Where("id IN (?)", assetIDs).Updates(map[string]interface{}{
			"maintainer_id": maintainerID,
			"state":         2,
		})
	}
	err := utils.DBError(result)
	return err
}

func (asset *assetDao) CheckAssetPropertyExist(assetID uint, key string) (bool, error) {
	var thisAsset *model.Asset
	result := db.Model(&model.Asset{}).Where("id = ?", assetID).
		First(&thisAsset, datatypes.JSONQuery("property").HasKey(key))

	if result.Error == gorm.ErrRecordNotFound {
		return false, nil
	} else if result.Error != nil {
		return false, utils.DBError(result)
	}

	return true, nil
}

func (asset *assetDao) SetAssetProperty(assetID uint, key string, value string) error {
	result := db.Model(&model.Asset{}).Where("id = ?", assetID).
		UpdateColumn("property", datatypes.JSONSet("property").Set(key, value))

	return utils.DBError(result)
}

func (asset *assetDao) GetAssetProperty(assetID uint) (*model.Asset, error) {
	var thisAsset *model.Asset

	result := db.Model(&model.Asset{}).Where("id = ?", assetID).First(&thisAsset)

	return thisAsset, utils.DBError(result)
}

func (asset *assetDao) GetAssetTask(assetID uint) ([]*model.Task, error) {
	var taskList []*model.Task
	var thisAsset *model.Asset

	result := db.Model(&model.Asset{}).Preload("TaskList.User").Preload("TaskList.Target").Where("id = ?", assetID).First(&thisAsset)

	taskList = thisAsset.TaskList

	return taskList, utils.DBError(result)
}

func (asset *assetDao) SearchDepartmentAsset(departmentID uint, req *define.SearchAssetReq) (assetList []*model.Asset, err error) {
	result := db.Model(&model.Asset{}).Where("department_id = ?", departmentID)

	if req.Name != "" {
		result = result.Where("name LIKE ?", req.Name)
	}

	if req.Description != "" {
		result = result.Where("description LIKE ?", req.Description)
	}

	if req.UserID != 0 {
		result = result.Where("user_id = ?", req.UserID)
	}

	if req.State < 5 {
		result = result.Where("state = ?", req.State)
	}

	if req.ClassID != 0 {
		result = result.Where("class_id = ?", req.ClassID)
	}

	if req.Key != "" {
		if req.Value == "" {
			result = result.Preload("Parent").Preload("User").
				Preload("Department").Preload("Class").Preload("Maintainer").
				Find(&assetList, datatypes.JSONQuery("property").HasKey(req.Key))
		} else {
			result = result.Preload("Parent").Preload("User").
				Preload("Department").Preload("Class").Preload("Maintainer").
				Find(&assetList, datatypes.JSONQuery("property").Equals(req.Value, req.Key))
		}
	} else {
		result = result.Preload("Parent").Preload("User").
			Preload("Department").Preload("Class").Preload("Maintainer").
			Find(&assetList)
	}

	if result.Error != nil {
		assetList = nil
		err = utils.DBError(result)
	} else {
		err = nil
	}

	return
}

func (asset *assetDao) GetDepartmentAssetCount(departmentID uint) (count int64, err error) {
	result := db.Model(&model.Asset{}).Where("department_id = ? and state <= ?", departmentID, 2).Count(&count)
	err = utils.DBError(result)
	return
}

func (asset *assetDao) GetDepartmentWarnAsset(departmentID uint) (assetList []*model.Asset, err error) {
	result := db.Model(&model.Asset{}).Where("department_id = ? and state <= ? and warn = ?", departmentID, 2, true).Find(&assetList)
	err = utils.DBError(result)
	return
}
