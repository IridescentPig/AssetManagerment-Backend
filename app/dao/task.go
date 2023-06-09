package dao

import (
	"asset-management/app/model"
	"asset-management/utils"
	"errors"
	"time"

	"gorm.io/gorm"
)

type taskDao struct {
}

var TaskDao *taskDao

func newTaskDao() *taskDao {
	return &taskDao{}
}

func init() {
	TaskDao = newTaskDao()
}

func (task *taskDao) Create(newTask model.Task) (uint, error) {
	result := db.Model(&model.Task{}).Create(&newTask)
	return newTask.ID, utils.DBError(result)
}

func (task *taskDao) Delete(id []uint) error {
	result := db.Model(&model.Task{}).Where("id in (?)", id).Delete(&model.Task{})
	return utils.DBError(result)
}

func (task *taskDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.Task{}).Where("id = ?", id).Updates(data)
	return utils.DBError(result)
}

func (task *taskDao) GetTaskByID(id uint) (*model.Task, error) {
	ret := &model.Task{}
	result := db.Model(&model.Task{}).Preload("User").Preload("Target").Preload("Department").Preload("AssetList.Class").Where("id = ?", id).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (task *taskDao) ModifyTaskType(id uint, TaskType uint) error {
	thisTask, err := task.GetTaskByID(id)
	if err != nil {
		return err
	}
	if thisTask == nil {
		err = errors.New("task doesn't exist")
		return err
	}
	thisTask.TaskType = TaskType
	err = utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&thisTask))
	return err
}

func (task *taskDao) ModifyTaskDescription(id uint, description string) error {
	thisTask, err := task.GetTaskByID(id)
	if err != nil {
		return err
	}
	if thisTask == nil {
		err = errors.New("task doesn't exist")
		return err
	}
	thisTask.TaskDescription = description
	err = utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&thisTask))
	return err
}

func (task *taskDao) GetTaskUser(id uint) (user model.User, err error) {
	thisTask, err := task.GetTaskByID(id)
	if err != nil {
		return
	}
	if thisTask == nil {
		err = errors.New("task doesn't exist")
		return
	}
	user = thisTask.User
	return
}

func (task *taskDao) ModifyTaskUser(id uint, user model.User) error {
	thisTask, err := task.GetTaskByID(id)
	if err != nil {
		return err
	}
	if thisTask == nil {
		err = errors.New("task doesn't exist")
		return err
	}
	thisTask.UserID = user.ID
	// thisTask.UserName = user.UserName
	err = utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&thisTask))
	return err
}

func (task *taskDao) GetTaskTarget(id uint) (user model.User, err error) {
	thisTask, err := task.GetTaskByID(id)
	if err != nil {
		return
	}
	if thisTask == nil {
		err = errors.New("task doesn't exist")
		return
	}
	user = thisTask.Target
	return
}

func (task *taskDao) ModifyTaskTarget(id uint, user model.User) error {
	thisTask, err := task.GetTaskByID(id)
	if err != nil {
		return err
	}
	if thisTask == nil {
		err = errors.New("task doesn't exist")
		return err
	}
	thisTask.TargetID = user.ID
	// thisTask.TargetName = user.UserName
	err = utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&thisTask))
	return err
}

// task and asset
func (task *taskDao) GetTaskAssetList(id uint) (list []*model.Asset, err error) {
	thisTask, err := task.GetTaskByID(id)
	if err != nil {
		return
	}
	if thisTask == nil {
		err = errors.New("task doesn't exist")
		return
	}
	list = thisTask.AssetList
	return
}

func (task *taskDao) ModifyAssetList(id uint, list []*model.Asset) error {
	thisTask, err := task.GetTaskByID(id)
	if err != nil {
		return err
	}
	if thisTask == nil {
		err = errors.New("task doesn't exist")
		return err
	}
	thisTask.AssetList = list
	err = utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&thisTask))
	return err
}

func (task *taskDao) GetTaskListByUserID(userID uint) (taskList []*model.Task, err error) {
	result := db.Model(&model.Task{}).Preload("User").
		Preload("Target").Preload("Department").
		Preload("AssetList.Class").Where("user_id = ?", userID).Find(&taskList)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

func (task *taskDao) GetTaskListByDepartmentID(departmentID uint) (taskList []*model.Task, err error) {
	result := db.Model(&model.Task{}).Preload("User").
		Preload("Target").Preload("Department").
		Preload("AssetList.Class").Where("department_id = ? and state <= ?", departmentID, 2).Find(&taskList)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	err = utils.DBError(result)
	return
}

func (task *taskDao) ModifyTaskState(taskID uint, state uint) error {
	result := db.Model(&model.Task{}).Where("id = ?", taskID).Updates(map[string]interface{}{
		"state":     state,
		"review_at": model.ModelTime(time.Now()),
	})
	return utils.DBError(result)
}
