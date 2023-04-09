package dao

import (
	"asset-management/app/model"
	"asset-management/utils"
	"errors"

	"github.com/shopspring/decimal"
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

func (asset *assetDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.Asset{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (asset *assetDao) AllUpdate(ids []int, data map[string]interface{}) error {
	result := db.Model(&model.Asset{}).Where("id IN (?)", ids).Updates(data)
	return utils.DBError(result)
}

func (asset *assetDao) Delete(id []uint) error {
	result := db.Model(&model.Asset{}).Where("id in (?)", id).Delete(&model.Asset{})
	return utils.DBError(result)
}

func (asset *assetDao) AllAsset() (list []model.Asset, err error) {
	result := db.Model(&model.Asset{}).Find(&list)
	err = utils.DBError(result)
	return
}

func (asset *assetDao) GetAssetByName(username string) (list []model.Asset, err error) {
	result := db.Model(&model.User{}).Where("name = ?", username).Find(&list)
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
		err = db.Model(&asset).Association("AssetClass").Find(&asset_class)
		if err != nil {
			return
		}
		asset.ClassID = asset_class.ID
	}
	err = utils.DBError(result)
	return
}

func (asset *assetDao) GetAssetByID(id int) (*model.Asset, error) {
	ret := &model.Asset{}
	result := db.Model(&model.Asset{}).Where("ID = ?", id).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	user := &model.User{}
	err := db.Model(&ret).Association("User").Find(&user)
	if err != nil {
		return nil, err
	}
	ret.UserID = user.ID
	asset_class := &model.AssetClass{}
	err = db.Model(&asset).Association("AssetClass").Find(&asset_class)
	if err != nil {
		return nil, err
	}
	ret.ClassID = asset_class.ID
	return ret, utils.DBError(result)
}

func (asset *assetDao) AssetCount() (count int64, err error) {
	result := db.Model(&model.Asset{}).Count(&count)
	err = utils.DBError(result)
	return
}

var asset_not_exist string = "asset doesn't exist"

func (asset *assetDao) ModifyAssetPrice(id int, price decimal.Decimal) error {
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

func (asset *assetDao) ModifyAssetDescription(id int, description string) error {
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

func (asset *assetDao) ModifyAssetPosition(id int, position string) error {
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

func (asset *assetDao) ExpireAsset(ids []int) error {
	return asset.AllUpdate(ids, map[string]interface{}{
		"expire": true,
		"price":  decimal.NewFromFloat(0),
	})
}

// asset and asset
func (asset *assetDao) GetSubAsset(id int) (assets []*model.Asset, err error) {
	query_asset, err := asset.GetAssetByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&query_asset).Where("parent_id = ?", query_asset.ID).Find(&assets))
	return
}

func (asset *assetDao) GetParentAsset(id int) (ParentAsset *model.Asset, err error) {
	query_asset, err := asset.GetAssetByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&query_asset).Where("id = ?", query_asset.ParentID).Find(&ParentAsset))
	return
}

func (asset *assetDao) ModifyParentAsset(ChildID int, ParentID int) error {
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
func (asset *assetDao) GetAssetUser(id int) (user model.User, err error) {
	query_asset, err := asset.GetAssetByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&model.Asset{}).Where("id = ?", query_asset.UserID).Find(&user))
	return
}

func (asset *assetDao) ModifyAssetUser(AssetID int, Username string) error {
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

// asset and asset_class
func (asset *assetDao) GetAssetClass(id int) (class model.AssetClass, err error) {
	query_asset, err := asset.GetAssetByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&model.AssetClass{}).Where("id = ?", query_asset.ClassID).Find(&class))
	return
}

var type_not_match = "type not match"

func (asset *assetDao) ModifyAssetClass(AssetID int, ClassID int) error {
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
	query_asset.UserID = target_class.ID
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&query_asset))
}
