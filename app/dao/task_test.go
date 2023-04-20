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
	assert.Equal(t, nil, err, err)
	err = TaskDao.Update(1, map[string]interface{}{
		"TaskType": 3,
	})
	assert.Equal(t, nil, err, err)

	_, err = TaskDao.GetTaskByID(1)
	assert.Equal(t, nil, err, err)

	err = TaskDao.ModifyTaskType(1, 1)
	assert.Equal(t, nil, err, err)

	err = TaskDao.ModifyTaskDescription(1, "suck")
	assert.Equal(t, nil, err, err)

	new_user, err := UserDao.GetUserByID(1)
	assert.Equal(t, nil, err, err)

	err = TaskDao.ModifyTaskUser(1, *new_user)
	assert.Equal(t, nil, err, err)

	_, err = TaskDao.GetTaskUser(1)
	assert.Equal(t, nil, err, "database error")

	err = TaskDao.ModifyTaskTarget(1, *new_user)
	assert.Equal(t, nil, err, "database error")

	_, err = TaskDao.GetTaskTarget(1)
	assert.Equal(t, nil, err, "database error")

	asset := model.Asset{
		Name: "a",
	}
	asset_list := []*model.Asset{&asset}

	err = TaskDao.ModifyAssetList(1, asset_list)
	assert.Equal(t, nil, err, "database error")

	_, err = TaskDao.GetTaskAssetList(1)
	assert.Equal(t, nil, err, "database error")

	err = TaskDao.Delete([]uint{1})
	assert.Equal(t, nil, err, err)

}
