package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

	"gorm.io/gorm"
)

type asyncDao struct {
}

var AsyncDao *asyncDao

func newAsyncDao() *asyncDao {
	return &asyncDao{}
}

func init() {
	AsyncDao = newAsyncDao()
}

func (asy *asyncDao) GetPendingTask() (task *model.AsyncTask, err error) {
	result := db.Model(&model.AsyncTask{}).Where("state = ?", 0).First(&task)
	if result.Error == gorm.ErrRecordNotFound {
		task = nil
		err = nil
	}
	err = utils.DBError(result)
	return
}
