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

func (assetclass *assetClassDao) UpdateByStruct(id uint, data model.AssetClass) error {
	result := db.Model(&model.AssetClass{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (assetclass *assetClassDao) AllUpdate(ids []uint, data map[string]interface{}) error {
	result := db.Model(&model.AssetClass{}).Where("id IN (?)", ids).Updates(data)
	return utils.DBError(result)
}

func (assetclass *assetClassDao) Delete(id []uint) error {
	result := db.Model(&model.AssetClass{}).Where("id in (?)", id).Delete(&model.AssetClass{})
	return utils.DBError(result)
}

func (assetclass *assetClassDao) GetAssetClassByID(id uint) (*model.AssetClass, error) {
	ret := &model.AssetClass{}
	result := db.Model(&model.AssetClass{}).Where("ID = ?", id).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	department := &model.Department{}
	err := db.Model(&ret).Association("Department").Find(&department)
	if err != nil {
		return nil, err
	}
	ret.Department = *department
	return ret, utils.DBError(result)
}

func (assetclass *assetClassDao) GetDepartmentDirectClass(departmentID uint) (assetClasses []*model.AssetClass, err error) {
	result := db.Model(&model.AssetClass{}).Where("department_id = ? and parent_id IS NULL", departmentID).Find(&assetClasses)
	if result.Error == gorm.ErrRecordNotFound {
		err = nil
	} else {
		err = utils.DBError(result)
	}
	return
}

// assetclass and assetclass
func (assetclass *assetClassDao) GetSubAssetClass(id uint) (assetClasses []*model.AssetClass, err error) {
	result := db.Model(&model.AssetClass{}).Where("parent_id = ?", id).Find(&assetClasses)
	if result.Error == gorm.ErrRecordNotFound {
		err = nil
	} else {
		err = utils.DBError(result)
	}
	return
}

func (assetclass *assetClassDao) GetParentAssetClass(id uint) (ParentAssetClass *model.AssetClass, err error) {
	query_asset, err := assetclass.GetAssetClassByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&query_asset).Where("id = ?", query_asset.ParentID).Find(&ParentAssetClass))
	return
}

func (assetclass *assetClassDao) ModifyParentAssetClass(ChildID uint, ParentID uint) error {
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
func (assetclass *assetClassDao) GetAssetClassDepartment(id uint) (department model.Department, err error) {
	query_asset, err := assetclass.GetAssetClassByID(id)
	if err != nil {
		return
	}
	department = query_asset.Department
	return
}

func (assetclass *assetClassDao) ModifyAssetClassDepartment(AssetClassID uint, DepartmentID uint) error {
	query_asset, err := assetclass.GetAssetClassByID(AssetClassID)
	if err != nil {
		return err
	}
	target_department, err := DepartmentDao.GetDepartmentByID(DepartmentID)
	if err != nil {
		return err
	}
	query_asset.DepartmentID = target_department.ID
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&query_asset))
}
