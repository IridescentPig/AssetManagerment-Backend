package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
)

type asyncService struct {
}

var AsyncService *asyncService

func newAsyncService() *asyncService {
	return &asyncService{}
}

func init() {
	AsyncService = newAsyncService()
}

func (asy *asyncService) GetUserAsyncTasks(userID uint) ([]*model.AsyncTask, error) {
	taskList, err := dao.AsyncDao.GetAsyncTaskListByUserID(userID)
	return taskList, err
}

func (asy *asyncService) CreateAsyncTask(userID uint, req *define.CreateAsyncTaskReq) error {
	task := model.AsyncTask{
		Type:      req.Type,
		UserID:    userID,
		ObjectKey: req.ObjectKey,
		FromTime:  req.FromTime,
		EntityID:  req.EntityID,
	}
	if req.DepartmentID != 0 {
		task.DepartmentID = req.DepartmentID
	}

	err := dao.AsyncDao.CreateAsyncTask(task)

	return err
}

func (asy *asyncService) GetAsyncTaskByID(taskID uint) (*model.AsyncTask, error) {
	return dao.AsyncDao.GetAsyncTaskByID(taskID)
}

func (asy *asyncService) ModifyAsyncTaskState(taskID uint, state uint) error {
	return dao.AsyncDao.ModifyAsyncTaskInfo(taskID, map[string]interface{}{
		"state": state,
	})
}
