package dao

import (
	"asset-management/app/model"
	"asset-management/utils"
	"log"

	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type mysqlHook struct {
}

var MysqlHook *mysqlHook

type logDao struct {
}

var LogDao *logDao

func newLogDao() *logDao {
	return &logDao{}
}

func init() {
	LogDao = newLogDao()
	MysqlHook = &mysqlHook{}
}

func (hook *mysqlHook) Fire(entry *logrus.Entry) error {
	_, err := entry.String()
	if err != nil {
		return err
	}
	go func() {
		mylog := &model.Log{
			Method:       entry.Data["method"].(string),
			URL:          entry.Data["url"].(string),
			Status:       entry.Data["status"].(int),
			ErrorCode:    entry.Data["error_code"].(int),
			ErrorMessage: entry.Data["error_message"].(string),
			UserID:       entry.Data["user_id"].(uint),
			Username:     entry.Data["username"].(string),
			EntityID:     entry.Data["entity_id"].(uint),
			DepartmentID: entry.Data["department_id"].(uint),
			Time:         (*model.ModelTime)(&entry.Time),
			Level:        entry.Level.String(),
			Message:      entry.Message,
		}

		result := db.Model(&model.Log{}).Create(mylog)
		if result.Error != nil {
			log.Println("Write log to database error: ", result.Error.Error())
		}
	}()
	// log := &model.Log{
	// 	Method:       entry.Data["method"].(string),
	// 	URL:          entry.Data["url"].(string),
	// 	Status:       entry.Data["status"].(int),
	// 	ErrorCode:    entry.Data["error_code"].(int),
	// 	ErrorMessage: entry.Data["error_message"].(string),
	// 	UserID:       entry.Data["user_id"].(uint),
	// 	Username:     entry.Data["username"].(string),
	// 	EntityID:     entry.Data["entity_id"].(uint),
	// 	DepartmentID: entry.Data["department_id"].(uint),
	// 	Time:         (*model.ModelTime)(&entry.Time),
	// 	Level:        entry.Level.String(),
	// 	Message:      entry.Message,
	// }

	// result := db.Model(&model.Log{}).Create(log)
	// return utils.DBError(result)
	return nil
}

func (hook *mysqlHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (mylog *logDao) GetLoginLogByEntityID(entityID uint, offset int, limit int) (logList []*model.Log, count int64, err error) {
	result := db.Model(&model.Log{}).Where("entity_id = ? and url = ?", entityID, "/user/login").Count(&count).Offset(offset).Limit(limit).Find(&logList)
	if result.Error == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	err = utils.DBError(result)
	return
}

func (mylog *logDao) GetLoginLogsForExport(entityID uint, fromTime *model.ModelTime, logType uint) (logList []*model.Log, err error) {
	var result *gorm.DB
	if logType == 0 {
		result = db.Model(&model.Log{}).Where("entity_id = ? and url = ? and time >= ?", entityID, "/user/login", fromTime)
	} else if logType == 1 {
		result = db.Model(&model.Log{}).Where("entity_id = ? and url = ? and time >= ? and status = ?", entityID, "/user/login", fromTime, 200)
	} else {
		result = db.Model(&model.Log{}).Where("entity_id = ? and url = ? and time >= ? and status <> ?", entityID, "/user/login", fromTime, 200)
	}

	result = result.Find(&logList)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

func (mylog *logDao) GetDataLogByEntityID(entityID uint, offset int, limit int) (logList []*model.Log, count int64, err error) {
	result := db.Model(&model.Log{}).Where("entity_id = ? and url <> ?", entityID, "/user/login").Count(&count).Offset(offset).Limit(limit).Find(&logList)
	if result.Error == gorm.ErrRecordNotFound {
		err = nil
		return
	}
	err = utils.DBError(result)
	return
}

func (mylog *logDao) GetDataLogsForExport(entityID uint, fromTime *model.ModelTime, logType uint) (logList []*model.Log, err error) {
	var result *gorm.DB
	if logType == 0 {
		result = db.Model(&model.Log{}).Where("entity_id = ? and url <> ? and time >= ?", entityID, "/user/login", fromTime)
	} else if logType == 1 {
		result = db.Model(&model.Log{}).Where("entity_id = ? and url <> ? and time >= ? and status = ?", entityID, "/user/login", fromTime, 200)
	} else {
		result = db.Model(&model.Log{}).Where("entity_id = ? and url <> ? and time >= ? and status <> ?", entityID, "/user/login", fromTime, 200)
	}

	result = result.Find(&logList)

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}
