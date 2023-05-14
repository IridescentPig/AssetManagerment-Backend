package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
)

type statService struct {
}

var StatService *statService

func newStatService() *statService {
	return &statService{}
}

func init() {
	StatService = newStatService()
}

func (stat *statService) GetDepartmentStat(departmentID uint) ([]*model.Stat, error) {
	return dao.StatDao.GetDepartmentStat(departmentID)
}

func (stat *statService) GetDepartmentAssetDistribution(departmentID uint) ([]*define.AssetDistribution, error) {
	return dao.StatDao.GetDepartmentAssetDistribution(departmentID)
}
