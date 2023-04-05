package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

	"gorm.io/gorm"
)

type entityDao struct {
}

var EntityDao *entityDao

func newEntityDao() *entityDao {
	return &entityDao{}
}

func init() {
	EntityDao = newEntityDao()
}

func (entity *entityDao) Create(newEntity model.Entity) error {
	result := db.Model(&model.Entity{}).Create(&newEntity)
	return utils.DBError(result)
}

func (entity *entityDao) Delete(id []uint) error {
	result := db.Model(&model.Entity{}).Where("id in (?)", id).Delete(&model.Entity{})
	return utils.DBError(result)
}

func (entity *entityDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.Entity{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (entity *entityDao) AllEntity() (list []model.Entity, err error) {
	result := db.Model(&model.Entity{}).Find(&list)
	err = utils.DBError(result)
	return
}

func (entity *entityDao) GetEntityByName(name string) (*model.Entity, error) {
	ret := &model.Entity{}
	result := db.Model(&model.Entity{}).Where("name = ?", name).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (entity *entityDao) GetEntitysByNames(name []string) (list []model.Entity, err error) {
	result := db.Model(&model.Entity{}).Where("name IN (?)", name).Order("id").Find(&list)
	err = utils.DBError(result)
	return
}

func (entity *entityDao) EntityCount() (count int64, err error) {
	result := db.Model(&model.Entity{}).Count(&count)
	err = utils.DBError(result)
	return
}

// entity and user
func (entity *entityDao) GetEntityAllUser(query_entity model.Entity) (users []*model.User, err error) {
	err = utils.DBError(db.Model(&query_entity).Where("id = ?", query_entity.ID).Preload("user").Find(&users))
	return
}

// entity and department
func (entity *entityDao) GetEntityAllDepartment(query_entity model.Entity) (departments []*model.Department, err error) {
	err = utils.DBError(db.Model(&query_entity).Where("id = ?", query_entity.ID).Preload("department").Find(&departments))
	return
}

func (entity *entityDao) GetEntitySubDepartment(query_entity model.Entity) (departments []*model.Department, err error) {
	err = utils.DBError(db.Model(&query_entity).Where("id = ?", query_entity.ID).Preload("department", "parent_id = 0 AND entity_id = ?", query_entity.ID).Find(&departments))
	return
}
