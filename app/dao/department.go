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

func (department *departmentDao) Update(id uint, data map[string]interface{}) error {
	result := db.Model(&model.Department{}).Where("id = ?", id).Updates(data)
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
	/*parent := &model.Department{}
	err := db.Model(&ret).Association("Department").Find(&parent)
	if err != nil {
		return nil, err
	}
	ret.Parent = parent*/
	entity := &model.Entity{}
	err := db.Model(&ret).Association("Entity").Find(&entity)
	if err != nil {
		return nil, err
	}
	ret.Entity = *entity
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (department *departmentDao) GetDepartmentSub(name string, entityID uint, departmentID uint) (*model.Department, error) {
	ret := &model.Department{}
	var result *gorm.DB
	if departmentID != 0 {
		result = db.Model(&model.Department{}).Preload("Parent").Preload("Entity").Where("name = ? and entity_id = ? and parent_id = ?", name, entityID, departmentID).First(ret)
	} else {
		result = db.Model(&model.Department{}).Preload("Parent").Preload("Entity").Where("name = ? and entity_id = ? and parent_id IS NULL", name, entityID).First(ret)
	}

	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return ret, utils.DBError(result)
}

func (department *departmentDao) GetDepartmentByID(id uint) (*model.Department, error) {
	ret := &model.Department{}
	result := db.Model(&model.Department{}).Preload("Parent").Preload("Entity").Where("id = ?", id).First(ret)
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

func (department *departmentDao) GetSubDepartmentByID(id uint) (departments []*model.Department, err error) {
	// query_department, err := department.GetDepartmentByID(id)
	// if err != nil {
	// 	return
	// }
	result := db.Model(&model.Department{}).Where("parent_id = ?", id).Find(&departments)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, nil
	}
	return departments, utils.DBError(result)
}

// department and department
func (department *departmentDao) GetSubDepartment(name string) (departments []*model.Department, err error) {
	query_department, err := department.GetDepartmentByName(name)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&query_department).Where("parent_id = ?", query_department.ID).Find(&departments))
	return
}

func (department *departmentDao) GetParentDepartment(name string) (departments *model.Department, err error) {
	query_department, err := department.GetDepartmentByName(name)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&query_department).Where("id = ?", query_department.ParentID).Find(&departments))
	return
}

func (department *departmentDao) ModifyParentDepartment(child_name string, parent_name string) error {
	child_department, err := department.GetDepartmentByName(child_name)
	if err != nil {
		return err
	}
	parent_department, err := department.GetDepartmentByName(parent_name)
	if err != nil {
		return err
	}
	child_department.ParentID = parent_department.ID
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&child_department))
}

// department and user
func (department *departmentDao) GetDepartmentDirectUser(name string) (users []*model.User, err error) {
	query_department, err := department.GetDepartmentByName(name)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&model.User{}).Where("department_id = ?", query_department.ID).Find(&users))
	return
}

func (department *departmentDao) GetDepartmentDirectUserByID(id uint) (users []*model.User, err error) {
	query_department, err := department.GetDepartmentByID(id)
	if err != nil {
		return
	}
	err = utils.DBError(db.Model(&model.User{}).Preload("Department").Preload("Entity").Where("department_id = ?", query_department.ID).Find(&users))
	return
}

func (department *departmentDao) GetDepartmentAllUserByID(id uint) (users []*model.User, err error) {
	direct_users, err := department.GetDepartmentDirectUserByID(id)
	if err != nil {
		return
	}
	users = append(users, direct_users...)
	departments, err := department.GetSubDepartmentByID(id)
	for _, dpm := range departments {
		indirect_users, in_err := department.GetDepartmentAllUserByID(dpm.ID)
		if in_err != nil {
			err = in_err
			return
		}
		users = append(users, indirect_users...)
	}
	return
}

func (department *departmentDao) GetDepartmentAllUser(name string) (users []*model.User, err error) {
	direct_users, err := department.GetDepartmentDirectUser(name)
	if err != nil {
		return
	}
	users = append(users, direct_users...)
	departments, err := department.GetSubDepartment(name)
	for _, dpm := range departments {
		indirect_users, in_err := department.GetDepartmentAllUser(dpm.Name)
		if in_err != nil {
			err = in_err
			return
		}
		users = append(users, indirect_users...)
	}
	return
}

// department and entity
func (department *departmentDao) GetDepartmentEntity(name string) (entity model.Entity, err error) {
	query_department, err := department.GetDepartmentByName(name)
	if err != nil {
		return
	}
	entity = query_department.Entity
	return
}

func (department *departmentDao) ModifyDepartmentEntity(department_name string, entity_name string) error {
	query_department, err := department.GetDepartmentByName(department_name)
	if err != nil {
		return err
	}
	target_entity, err := EntityDao.GetEntityByName(entity_name)
	if err != nil {
		return err
	}
	query_department.Entity = *target_entity
	return utils.DBError(db.Session(&gorm.Session{FullSaveAssociations: true}).Updates(&query_department))
}

func (department *departmentDao) GetDepartmentManager(id uint) (managers []*model.User, err error) {
	err = utils.DBError(db.Model(&model.User{}).Where("department_id = ? and department_super = ?", id, true).Find(&managers))
	return
}
