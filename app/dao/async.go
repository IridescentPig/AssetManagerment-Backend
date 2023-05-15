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

func (asy *asyncDao) GetAsyncTaskListByUserID(userID uint) (taskList []*model.AsyncTask, err error) {
	result := db.Model(&model.AsyncTask{}).Where("user_id = ?", userID).Find(&taskList)
	err = utils.DBError(result)
	return
}

func (asy *asyncDao) CreateAsyncTask(newTask model.AsyncTask) (err error) {
	result := db.Model(&model.AsyncTask{}).Preload("User").Create(&newTask)
	err = utils.DBError(result)
	return
}

func (asy *asyncDao) ModifyAsyncTaskInfo(taskID uint, data map[string]interface{}) (err error) {
	result := db.Model(&model.AsyncTask{}).Updates(data)
	err = utils.DBError(result)
	return
}
