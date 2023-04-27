package service

import (
	"asset-management/app/dao"
	"asset-management/app/model"
)

type logService struct {
}

var LogService *logService

func newLogService() *logService {
	return &logService{}
}

func init() {
	LogService = newLogService()
}

func (mylog *logService) GetLoginLog(entityID uint) (logList []*model.Log, err error) {
	logList, err = dao.LogDao.GetLoginLogByEntityID(entityID)
	return
}

func (mylog *logService) GetDataLog(entityID uint) (logList []*model.Log, err error) {
	logList, err = dao.LogDao.GetDataLogByEntityID(entityID)
	return
}
