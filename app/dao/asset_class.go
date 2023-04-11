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

// assetclass and assetclass
func (assetclass *assetClassDao) GetSubAssetClass(id int) (assetClasses []*model.AssetClass, err error) {
	query_asset, err := assetclass.GetAssetClassByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&query_asset).Where("parent_id = ?", query_asset.ID).Find(&assetClasses))
	return
}

func (assetclass *assetClassDao) GetParentAssetClass(id int) (ParentAssetClass *model.AssetClass, err error) {
	query_asset, err := assetclass.GetAssetClassByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&query_asset).Where("id = ?", query_asset.ParentID).Find(&ParentAssetClass))
	return
}

func (assetclass *assetClassDao) ModifyParentAssetClass(ChildID int, ParentID int) error {
	child_asset, err := assetclass.GetAssetClassByID(ChildID)
	if err != nil {
		return err
	}
	parent_asset, err := assetclass.GetAssetClassByID(ParentID)
	if err != nil {
		return err
	}
	child_asset.ParentID = parent_asset.ID
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&child_asset))
}

// assetclass and entity
func (assetclass *assetClassDao) GetAssetClassEntity(id int) (entity model.Entity, err error) {
	query_department, err := assetclass.GetAssetClassByID(id)
	if err != nil {
		return
	}
	entity = query_department.Entity
	return
}

func (assetclass *assetClassDao) ModifyDepartmentEntity(AssetClassID int, EntityID int) error {
	query_department, err := assetclass.GetAssetClassByID(AssetClassID)
	if err != nil {
		return err
	}
	target_entity, err := EntityDao.GetEntityByID(EntityID)
	if err != nil {
		return err
	}
	query_department.Entity = *target_entity
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&query_department))
}
