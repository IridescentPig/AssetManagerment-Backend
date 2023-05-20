package dao

import (
	"testing"
)

func TestLog(t *testing.T) {
	Init()

	MysqlHook.Levels()
	LogDao.GetLoginLogByEntityID(1, -1, -1)
	LogDao.GetLoginLogByEntityID(9, -1, -1)
	LogDao.GetDataLogByEntityID(1, -1, -1)
	LogDao.GetDataLogByEntityID(9, -1, -1)
}
