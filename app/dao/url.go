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

func (url *urlDao) Delete(names []string) error {
	result := db.Model(&model.Url{}).Where("name in (?)", names).Delete(&model.Url{})
	return utils.DBError(result)
}

func (url *urlDao) Update(name string, data map[string]interface{}) error {
	result := db.Model(&model.Url{}).Where("name = ?", name).Updates(data)
	return utils.DBError(result)
}

func (url *urlDao) GetUrlByName(name string) (*model.Url, error) {
	ret := &model.Url{}
	result := db.Model(&model.Url{}).Preload("Entity").Where("name = ?", name).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}
