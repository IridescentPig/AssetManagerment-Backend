package service

import (
	"asset-management/app/dao"
	"asset-management/app/define"
	"asset-management/app/model"
)

type taskService struct {
}

var TaskService *taskService

func newTaskService() *taskService {
	return &taskService{}
}

func init() {
	TaskService = newTaskService()
}

func (task *taskService) CreateTask(req define.CreateTaskReq, userID uint, departmentID uint, assetList []*model.Asset) error {
	err := dao.TaskDao.Create(model.Task{
		TaskType:        req.TaskType,
		TaskDescription: req.TaskDescription,
		UserID:          userID,
		DepartmentID:    departmentID,
		TargetID:        req.TargetID,
		AssetList:       assetList,
	})

	return err
}
