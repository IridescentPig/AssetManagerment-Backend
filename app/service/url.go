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

func (url *urlService) CreateUrl(req define.CreateOrModifyUrlReq, entity_id uint) error {
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

func (url *urlService) ModifyUrlInfo(req define.CreateOrModifyUrlReq, entity_id uint) error {
	err := dao.UrlDao.Update(req.Name, entity_id, map[string]interface{}{
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

func (url *urlService) GetUrlsByEntity(entity_id uint) ([]*model.Url, error) {
	return dao.UrlDao.GetUrlsByEntity(entity_id)
}
