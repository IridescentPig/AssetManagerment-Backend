package dao

import (
	"asset-management/app/model"
	"asset-management/utils"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type asyncDao struct {
}

var AsyncDao *asyncDao
var newLogger logger.Interface

func newAsyncDao() *asyncDao {
	return &asyncDao{}
}

func init() {
	AsyncDao = newAsyncDao()
	newLogger = logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold:             time.Second,   // Slow SQL threshold
			LogLevel:                  logger.Silent, // Log level
			IgnoreRecordNotFoundError: true,          // Ignore ErrRecordNotFound error for logger
			ParameterizedQueries:      true,          // Don't include params in the SQL log
			Colorful:                  true,          // Disable color
		},
	)
}

func (asy *asyncDao) GetPendingTask() (task *model.AsyncTask, err error) {
	result := db.Session(&gorm.Session{Logger: newLogger}).Model(&model.AsyncTask{}).Where("state = ?", 0).First(&task)
	if result.Error == gorm.ErrRecordNotFound {
		task = nil
		err = nil
	} else {
		err = utils.DBError(result)
	}
	return
}

func (asy *asyncDao) GetAsyncTaskListByUserID(userID uint) (taskList []*model.AsyncTask, err error) {
	result := db.Model(&model.AsyncTask{}).Preload("User").Where("user_id = ?", userID).Find(&taskList)
	err = utils.DBError(result)
	return
}

func (asy *asyncDao) CreateAsyncTask(newTask model.AsyncTask) (err error) {
	result := db.Model(&model.AsyncTask{}).Create(&newTask)
	err = utils.DBError(result)
	return
}

func (asy *asyncDao) ModifyAsyncTaskInfo(taskID uint, data map[string]interface{}) (err error) {
	result := db.Model(&model.AsyncTask{}).Where("id = ?", taskID).Updates(data)
	err = utils.DBError(result)
	return
}

func (asy *asyncDao) GetAsyncTaskByID(taskID uint) (task *model.AsyncTask, err error) {
	result := db.Model(&model.AsyncTask{}).Preload("User").Where("id = ?", taskID).First(&task)
	if result.Error == gorm.ErrRecordNotFound {
		err = nil
	} else {
		err = utils.DBError(result)
	}

	return
}
