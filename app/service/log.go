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

func (mylog *logService) GetLoginLog(entityID uint, page_size uint, page_num uint) (logList []*model.Log, count int64, err error) {
	offset := page_size * page_num
	limit := page_size
	logList, count, err = dao.LogDao.GetLoginLogByEntityID(entityID, int(offset), int(limit))
	return
}

func (mylog *logService) GetDataLog(entityID uint, page_size uint, page_num uint) (logList []*model.Log, count int64, err error) {
	offset := page_size * page_num
	limit := page_size
	logList, count, err = dao.LogDao.GetDataLogByEntityID(entityID, int(offset), int(limit))
	return
}
