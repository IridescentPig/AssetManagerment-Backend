package dao

import (
	"asset-management/app/model"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTask(t *testing.T) {
	Init()

	user := model.User{
		UserName: "test",
		Password: "123456",
	}
	err := UserDao.Create(user)
	assert.Equal(t, nil, err, "database error")

	user = model.User{
		UserName: "admin",
		Password: "123456",
	}
	err = UserDao.Create(user)
	assert.Equal(t, nil, err, err)

	task := model.Task{
		TaskType:        2,
		TaskDescription: "aaa",
	}
	err = TaskDao.Create(task)
	task = model.Task{
		TaskType:        3,
		TaskDescription: "asfdsag",
	}
	err = TaskDao.Create(task)
	task = model.Task{
		TaskType:        1,
		TaskDescription: "asgrah",
	}
	err = TaskDao.Create(task)
	assert.Equal(t, nil, err, err)
	err = TaskDao.Update(1, map[string]interface{}{
		"TaskType": 3,
	})
	err = TaskDao.Update(2, map[string]interface{}{
		"TaskType": 1,
	})
	err = TaskDao.Update(3, map[string]interface{}{
		"TaskType": 2,
	})
	assert.Equal(t, nil, err, err)

	_, err = TaskDao.GetTaskByID(1)
	assert.Equal(t, nil, err, err)
	_, err = TaskDao.GetTaskByID(2)
	_, err = TaskDao.GetTaskByID(9)

	err = TaskDao.ModifyTaskType(1, 1)
	assert.Equal(t, nil, err, err)
	err = TaskDao.ModifyTaskType(3, 2)
	err = TaskDao.ModifyTaskType(8, 3)

	err = TaskDao.ModifyTaskDescription(1, "suck")
	assert.Equal(t, nil, err, err)
	err = TaskDao.ModifyTaskDescription(3, "s2weck")
	err = TaskDao.ModifyTaskDescription(8, "suck")

	new_user, err := UserDao.GetUserByID(1)
	assert.Equal(t, nil, err, err)

	err = TaskDao.ModifyTaskUser(1, *new_user)
	assert.Equal(t, nil, err, err)
	err = TaskDao.ModifyTaskUser(2, *new_user)
	err = TaskDao.ModifyTaskUser(9, *new_user)

	_, err = TaskDao.GetTaskUser(1)
	assert.Equal(t, nil, err, "database error")
	_, err = TaskDao.GetTaskUser(3)
	_, err = TaskDao.GetTaskUser(9)

	err = TaskDao.ModifyTaskTarget(1, *new_user)
	assert.Equal(t, nil, err, "database error")
	err = TaskDao.ModifyTaskTarget(2, *new_user)
	err = TaskDao.ModifyTaskTarget(9, *new_user)

	_, err = TaskDao.GetTaskTarget(1)
	assert.Equal(t, nil, err, "database error")
	_, err = TaskDao.GetTaskTarget(2)
	_, err = TaskDao.GetTaskTarget(9)

	asset := model.Asset{
		Name: "a",
	}
	asset_list := []*model.Asset{&asset}

	err = TaskDao.ModifyAssetList(1, asset_list)
	assert.Equal(t, nil, err, "database error")
	err = TaskDao.ModifyAssetList(2, asset_list)
	err = TaskDao.ModifyAssetList(9, asset_list)

	_, err = TaskDao.GetTaskAssetList(1)
	assert.Equal(t, nil, err, "database error")
	_, err = TaskDao.GetTaskAssetList(2)
	_, err = TaskDao.GetTaskAssetList(9)

	err = TaskDao.Delete([]uint{1})
	assert.Equal(t, nil, err, err)
	err = TaskDao.Delete([]uint{2})
	err = TaskDao.Delete([]uint{7, 8, 9})

	TaskDao.GetTaskListByUserID(1)
	TaskDao.GetTaskListByUserID(2)
	TaskDao.GetTaskListByDepartmentID(1)
	TaskDao.GetTaskListByDepartmentID(2)
	TaskDao.ModifyTaskState(1, 2)
	TaskDao.ModifyTaskState(1, 8)
	TaskDao.ModifyTaskState(4, 2)

}
