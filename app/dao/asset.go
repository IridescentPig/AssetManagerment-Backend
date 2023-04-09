package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

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

func (asset *assetDao) Delete(id []uint) error {
	result := db.Model(&model.User{}).Where("id in (?)", id).Delete(&model.User{})
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
	ret.User = *user
	return ret, utils.DBError(result)
}

func (asset *assetDao) AssetCount() (count int64, err error) {
	result := db.Model(&model.Asset{}).Count(&count)
	err = utils.DBError(result)
	return
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

func (asset *assetDao) ModifyDepartUser(AssetID int, Username string) error {
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
