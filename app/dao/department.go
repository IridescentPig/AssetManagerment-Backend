package dao

type departmentDao struct {
}

var DepartmentDao *departmentDao

func newDepartmentDao() *departmentDao {
	return &departmentDao{}
}

func init() {
	DepartmentDao = newDepartmentDao()
}
