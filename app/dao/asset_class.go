package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

	"gorm.io/gorm"
)

type assetClassDao struct {
}

var AssetClassDao *assetClassDao

func newAssetClassDao() *assetClassDao {
	return &assetClassDao{}
}

func init() {
	AssetClassDao = newAssetClassDao()
}

func (assetclass *assetClassDao) Create(newAssetClass model.AssetClass) error {
	result := db.Model(&model.AssetClass{}).Create(&newAssetClass)
	return utils.DBError(result)
}

func (assetclass *assetClassDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.AssetClass{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (assetclass *assetClassDao) AllUpdate(ids []int, data map[string]interface{}) error {
	result := db.Model(&model.AssetClass{}).Where("id IN (?)", ids).Updates(data)
	return utils.DBError(result)
}

func (assetclass *assetClassDao) Delete(id []uint) error {
	result := db.Model(&model.AssetClass{}).Where("id in (?)", id).Delete(&model.AssetClass{})
	return utils.DBError(result)
}

func (assetclass *assetClassDao) GetAssetClassByID(id int) (*model.AssetClass, error) {
	ret := &model.AssetClass{}
	result := db.Model(&model.AssetClass{}).Where("ID = ?", id).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	entity := &model.Entity{}
	err := db.Model(&ret).Association("Entity").Find(&entity)
	if err != nil {
		return nil, err
	}
	ret.Entity = *entity
	return ret, utils.DBError(result)
}
