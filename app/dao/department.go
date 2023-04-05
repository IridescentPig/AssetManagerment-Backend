package dao

import (
	"asset-management/app/model"
	"asset-management/utils"

	"gorm.io/gorm"
)

type departmentDao struct {
}

var DepartmentDao *departmentDao

func newDepartmentDao() *departmentDao {
	return &departmentDao{}
}

func init() {
	DepartmentDao = newDepartmentDao()
}

func (department *departmentDao) Create(newDepartment model.Department) error {
	result := db.Model(&model.Department{}).Create(&newDepartment)
	return utils.DBError(result)
}

func (department *departmentDao) Delete(id []uint) error {
	result := db.Model(&model.Department{}).Where("id in (?)", id).Delete(&model.Department{})
	return utils.DBError(result)
}

func (department *departmentDao) AllDepartment() (list []model.Department, err error) {
	result := db.Model(&model.Department{}).Find(&list)
	err = utils.DBError(result)
	return
}

func (department *departmentDao) GetDepartmentByName(name string) (*model.Department, error) {
	ret := &model.Department{}
	result := db.Model(&model.Department{}).Where("name = ?", name).First(ret)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (department *departmentDao) GetDepartmentsByNames(name []string) (list []model.Department, err error) {
	result := db.Model(&model.Department{}).Where("name IN (?)", name).Order("id").Find(&list)
	err = utils.DBError(result)
	return
}

func (department *departmentDao) DepartmentCount() (count int64, err error) {
	result := db.Model(&model.Entity{}).Count(&count)
	err = utils.DBError(result)
	return
}

// department and department
func (department *departmentDao) GetSubDepartment(query_department model.Department) (departments []*model.Department, err error) {
	err = utils.DBError(db.Model(&query_department).Where("parent_id = ?", query_department.ID).Find(&departments))
	return
}

func (department *departmentDao) GetParentDepartment(query_department model.Department) (departments *model.Department, err error) {
	err = utils.DBError(db.Model(&query_department).Where("id = ?", query_department.ParentID).Find(&departments))
	return
}

// department and user
func (department *departmentDao) GetDepartmentDirectUser(query_department model.Department) (users []*model.User, err error) {
	err = utils.DBError(db.Model(&query_department).Where("ID = ?", query_department.ID).Preload("user").Find(&users))
	return
}

func (department *departmentDao) GetDepartmentAllUser(query_department model.Department) (users []*model.User, err error) {
	direct_users, err := department.GetDepartmentDirectUser(query_department)
	if err != nil {
		return
	}
	users = append(users, direct_users...)
	departments, err := department.GetSubDepartment(query_department)
	for _, dpm := range departments {
		indirect_users, in_err := department.GetDepartmentAllUser(*dpm)
		if in_err != nil {
			err = in_err
			return
		}
		users = append(users, indirect_users...)
	}
	return
}
