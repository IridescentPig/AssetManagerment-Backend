package service

import (
	"asset-management/app/dao"
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
