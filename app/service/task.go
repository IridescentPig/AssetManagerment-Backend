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

func (task *taskService) GetTasksByUserID(userID uint) (taskList []*model.Task, err error) {
	taskList, err = dao.TaskDao.GetTaskListByUserID(userID)
	return
}

func (task *taskService) GetTasksByDepartmentID(departmentID uint) (taskList []*model.Task, err error) {
	taskList, err = dao.TaskDao.GetTaskListByDepartmentID(departmentID)
	return
}

func (task *taskService) GetTaskInfoByID(taskID uint) (taskInfo *model.Task, err error) {
	taskInfo, err = dao.TaskDao.GetTaskByID(taskID)
	return
}

func (task *taskService) ModifyTaskState(taskID uint, state uint) error {
	err := dao.TaskDao.ModifyTaskState(taskID, state)
	return err
}
