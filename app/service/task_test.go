package service

import (
	"asset-management/app/define"
	"asset-management/app/model"
	"testing"
)

func TestTask(t *testing.T) {
	InitForTest()

	req := define.CreateTaskReq{
		TaskType:        1,
		TaskDescription: "kkk",
		TargetID:        1,
	}

	TaskService.CreateTask(req, 1, 1, []*model.Asset{})

	TaskService.GetTasksByUserID(1)
	TaskService.GetTasksByDepartmentID(1)
	TaskService.GetTaskInfoByID(1)
	TaskService.ModifyTaskState(1, 2)
}
