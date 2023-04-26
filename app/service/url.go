package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
)

type urlService struct {
}

var UrlService *urlService

func newUrlService() *urlService {
	return &urlService{}
}

func init() {
	UrlService = newUrlService()
}

func (url *urlService) CreateUrl(req define.CreateUrlReq, entity_id uint) error {
	err := dao.UrlDao.Create(model.Url{
		Name:            req.Name,
		Url:             req.Url,
		EntityID:        entity_id,
		DepartmentSuper: req.DepartmentSuper,
		EntitySuper:     req.EntitySuper,
		SystemSuper:     req.SystemSuper,
	})
	return err
}

func (url *urlService) ModifyUrlInfo(req define.ModifyUrlReq, entity_id uint) error {
	err := dao.UrlDao.Update(req.OldName, entity_id, map[string]interface{}{
		"name":             req.Name,
		"url":              req.Url,
		"department_super": req.DepartmentSuper,
		"entity_super":     req.EntitySuper,
		"system_super":     req.SystemSuper,
	})
	return err
}

func (url *urlService) DeleteUrl(name string, entity_id uint) error {
	err := dao.UrlDao.Delete([]string{name}, entity_id)
	return err
}

func (url *urlService) GetUrlsByEntity(entity_id uint, DepartmentSuper bool, EntitySuper bool, SystemSuper bool) ([]*model.Url, error) {
	return dao.UrlDao.GetUrlsByEntity(entity_id, DepartmentSuper, EntitySuper, SystemSuper)
}

func (url *urlService) CheckIfUrlExists(name string, entity_id uint) (exists bool, err error) {
	target_url, err := dao.UrlDao.GetUrlByName(name, entity_id)
	if target_url == nil {
		return false, err
	}
	return true, err
}
