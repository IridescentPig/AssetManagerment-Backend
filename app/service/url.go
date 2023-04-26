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
		Name:     req.Name,
		Url:      req.Url,
		EntityID: entity_id,
	})
	return err
}

func (url *urlService) ModifyUrlInfo(req define.CreateOrModifyUrlReq) error {
	err := dao.UrlDao.Update(req.Name, map[string]interface{}{
		"name": req.Name,
		"url":  req.Url,
	})
	return err
}

func (url *urlService) DeleteUrl(name string) error {
	err := dao.UrlDao.Delete([]string{name})
	return err
}
