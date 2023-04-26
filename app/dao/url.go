package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

	"gorm.io/gorm"
)

type urlDao struct {
}

var UrlDao *urlDao

func newUrlDao() *urlDao {
	return &urlDao{}
}

func init() {
	UrlDao = newUrlDao()
}

func (url *urlDao) Create(newUrl model.Url) error {
	result := db.Model(&model.Url{}).Create(&newUrl)
	return utils.DBError(result)
}

func (url *urlDao) Delete(names []string, entity_id uint) error {
	result := db.Model(&model.Url{}).Where("name in (?) and entity_id = ?", names, entity_id).Delete(&model.Url{})
	return utils.DBError(result)
}

func (url *urlDao) Update(name string, entity_id uint, data map[string]interface{}) error {
	result := db.Model(&model.Url{}).Where("name = ? and entity_id = ?", name, entity_id).Updates(data)
	return utils.DBError(result)
}

func (url *urlDao) GetUrlByName(name string, entity_id uint) (*model.Url, error) {
	ret := &model.Url{}
	result := db.Model(&model.Url{}).Preload("Entity").Where("name = ? and entity_id = ?", name, entity_id).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (url *urlDao) GetUrlsByEntity(entity_id uint, DepartmentSuper bool, EntitySuper bool, SystemSuper bool) (urls []*model.Url, err error) {
	all_urls := db.Model(&model.Url{}).Preload("Entity").
		Where("entity_id = ?", entity_id)
	if !DepartmentSuper {
		all_urls = all_urls.Not("department_super = true")
	}
	if !EntitySuper {
		all_urls = all_urls.Not("entity_super = true")
	}
	if !SystemSuper {
		all_urls = all_urls.Not("system_super = true")
	}
	err = utils.DBError(all_urls.Find(&urls))
	return
}
