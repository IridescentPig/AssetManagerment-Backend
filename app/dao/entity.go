package dao

import "asset-management/utils"

type Entity struct {
	ID   uint   `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name string `gorm:"column:name" json:"name"`
}

type entityDao struct {
}

var EntityDao *entityDao

func newEntityDao() *entityDao {
	return &entityDao{}
}

func init() {
	EntityDao = newEntityDao()
}

func (entity *entityDao) Create(entityname string) error {
	newEntity := Entity{Name: entityname}
	result := db.Model(&Entity{}).Create(&newEntity)
	return utils.DB_error(result)
}

func (entity *entityDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&Entity{}).Where("id = ?", id).Updates(data)
	return utils.DB_error(result)
}

func (entity *entityDao) Delete(id []uint) error {
	result := db.Model(&Entity{}).Where("id in (?)", id).Delete(&Entity{})
	return utils.DB_error(result)
}

func (entity *entityDao) AllEntity() (list []Entity, err error) {
	result := db.Model(&Entity{}).Find(&list)
	err = utils.DB_error(result)
	return
}

func (entity *entityDao) AllEntityWhere(query interface{}, args ...interface{}) (list []Entity, err error) {
	result := db.Model(&Entity{}).Where(query, args...).Find(&list)
	err = utils.DB_error(result)
	return
}

func (entity *entityDao) OneEntityWhere(query interface{}, args ...interface{}) (record Entity, err error) {
	result := db.Model(&Entity{}).Where(query, args...).First(&record)
	err = utils.DB_error(result)
	return
}

func (entity *entityDao) EntityCount() (count int64, err error) {
	result := db.Model(&Entity{}).Count(&count)
	err = utils.DB_error(result)
	return
}

func (entity *entityDao) EntityCountWhere(query interface{}, args ...interface{}) (count int64, err error) {
	result := db.Model(&Entity{}).Where(query, args...).Count(&count)
	err = utils.DB_error(result)
	return
}
