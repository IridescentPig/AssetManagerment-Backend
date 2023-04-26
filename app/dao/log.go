package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

	"github.com/sirupsen/logrus"
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

	log := &model.Log{
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

	result := db.Model(&model.Log{}).Create(log)
	return utils.DBError(result)
}

func (hook *mysqlHook) Levels() []logrus.Level {
	return logrus.AllLevels
}
